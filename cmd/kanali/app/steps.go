package app

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	pluginPkg "plugin"
	"strconv"
	"strings"
	"time"

	"github.com/northwesternmutual/kanali/cmd/kanali/app/options"
	"github.com/northwesternmutual/kanali/pkg/apis/kanali.io/v2"
	kanaliErrors "github.com/northwesternmutual/kanali/pkg/errors"
	"github.com/northwesternmutual/kanali/pkg/logging"
	"github.com/northwesternmutual/kanali/pkg/metrics"
	plugin "github.com/northwesternmutual/kanali/pkg/plugin"
	store "github.com/northwesternmutual/kanali/pkg/store"
	tags "github.com/northwesternmutual/kanali/pkg/tags"
	"github.com/northwesternmutual/kanali/pkg/utils"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"k8s.io/api/core/v1"
	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/informers/core"
)

const (
	pluginSymbolName = "Plugin"
)

type flow []step

type step interface {
	getName() string
	do(ctx context.Context, proxy *v2.ApiProxy, k8sCoreClient core.Interface, metrics *metrics.Metrics, w http.ResponseWriter, r *http.Request, resp *http.Response, trace opentracing.Span) error
}

// httpClient allows for mocking an http client to assist in testing
type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func (f *flow) add(steps ...step) {
	for _, s := range steps {
		*f = append(*f, s)
	}
}

func (f *flow) play(ctx context.Context, proxy *v2.ApiProxy, k8sCoreClient core.Interface, metrics *metrics.Metrics, w http.ResponseWriter, r *http.Request, resp *http.Response, trace opentracing.Span) error {
	logger := logging.WithContext(ctx)

	for _, step := range *f {
		logger.Debug(fmt.Sprintf("playing step %s", step.getName()))
		err := step.do(ctx, proxy, k8sCoreClient, metrics, w, r, resp, trace)
		if err == nil {
			continue
		}
		trace.SetTag(tags.Error, true)
		trace.LogKV(
			"event", tags.Error,
			"error.message", err.Error(),
		)
		return err
	}
	return nil
}

type mockTargetStep struct{}

func (step mockTargetStep) getName() string {
	return "mock target"
}

func (step mockTargetStep) do(ctx context.Context, proxy *v2.ApiProxy, k8sCoreClient core.Interface, m *metrics.Metrics, w http.ResponseWriter, r *http.Request, resp *http.Response, trace opentracing.Span) error {

	if !mockTargetDefined(proxy) || !mockTargetEnabled(proxy) {
		return nil
	}

	logger := logging.WithContext(ctx)

	targetPath := utils.ComputeTargetPath(proxy.Spec.Source.Path, proxy.Spec.Target.Path, utils.ComputeURLPath(r.URL))

	untypedMr, err := store.MockTargetStore.Get(proxy.ObjectMeta.Namespace, proxy.Spec.Target.Backend.Mock.MockTargetName, targetPath, r.Method)
	if err != nil {
		logger.Error(err.Error())
		return &kanaliErrors.StatusError{Code: http.StatusInternalServerError, Err: fmt.Errorf("error retrieving mock target")}
	}
	if untypedMr == nil {
		return &kanaliErrors.StatusError{Code: http.StatusNotFound, Err: errors.New("mock target not found")}
	}
	mr, ok := untypedMr.(v2.Route)
	if !ok {
		logger.Error("mock target response not expected type")
		return &kanaliErrors.StatusError{Code: http.StatusNotFound, Err: errors.New("mock target not found")}
	}

	upstreamHeaders := http.Header{}
	for k, v := range mr.Headers {
		upstreamHeaders.Add(k, v)
	}

	// create a fake response
	responseRecorder := &httptest.ResponseRecorder{
		Code:      mr.StatusCode,
		Body:      bytes.NewBuffer(mr.Body),
		HeaderMap: upstreamHeaders,
	}

	m.Add(metrics.Metric{Name: "http_response_code", Value: strconv.Itoa(mr.StatusCode), Index: true})

	*resp = *responseRecorder.Result()

	return nil

}

func mockTargetDefined(proxy *v2.ApiProxy) bool {
	return len(proxy.Spec.Target.Backend.Mock.MockTargetName) > 0
}

func mockTargetEnabled(proxy *v2.ApiProxy) bool {
	return viper.GetBool(options.FlagProxyEnableMockResponses.GetLong())
}

type pluginsOnRequestStep struct{}

func (step pluginsOnRequestStep) getName() string {
	return "plugin onrequest"
}

func (step pluginsOnRequestStep) do(ctx context.Context, proxy *v2.ApiProxy, k8sCoreClient core.Interface, m *metrics.Metrics, w http.ResponseWriter, r *http.Request, resp *http.Response, trace opentracing.Span) error {

	for _, plugin := range proxy.Spec.Plugins {
		p, err := getPlugin(plugin)
		if err != nil {
			return err
		}
		if err := doOnRequest(ctx, m, plugin, *proxy, r, trace, *p); err != nil {
			return err
		}
	}
	return nil
}

func doOnRequest(ctx context.Context, m *metrics.Metrics, plugin v2.Plugin, proxy v2.ApiProxy, req *http.Request, span opentracing.Span, p plugin.Plugin) (e error) {
	logger := logging.WithContext(ctx)

	defer func() {
		if r := recover(); r != nil {
			logger.Error(fmt.Sprintf("OnRequest paniced: %v", r))
			e = errors.New("OnRequest paniced")
		}
	}()

	sp := opentracing.StartSpan(fmt.Sprintf("PLUGIN: ON_REQUEST: %s", plugin.Name), opentracing.ChildOf(span.Context()))
	defer sp.Finish()

	return p.OnRequest(ctx, plugin.Config, m, proxy, req, sp)
}

// PluginsOnResponseStep is factory that defines a step responsible for
// executing the on response lifecycle hook for all the defined plugins
type pluginsOnResponseStep struct{}

// GetName retruns the name of the PluginsOnResponseStep step
func (step pluginsOnResponseStep) getName() string {
	return "Plugin OnResponse"
}

// Do executes the logic of the PluginsOnResponseStep step
func (step pluginsOnResponseStep) do(ctx context.Context, proxy *v2.ApiProxy, k8sCoreClient core.Interface, m *metrics.Metrics, w http.ResponseWriter, r *http.Request, resp *http.Response, trace opentracing.Span) error {

	for _, plugin := range proxy.Spec.Plugins {
		p, err := getPlugin(plugin)
		if err != nil {
			return err
		}
		if err := doOnResponse(ctx, m, plugin, *proxy, r, resp, trace, *p); err != nil {
			return err
		}
	}

	return nil
}

func doOnResponse(ctx context.Context, m *metrics.Metrics, plugin v2.Plugin, proxy v2.ApiProxy, req *http.Request, resp *http.Response, span opentracing.Span, p plugin.Plugin) (e error) {
	logger := logging.WithContext(ctx)

	defer func() {
		if r := recover(); r != nil {
			logger.Error(fmt.Sprintf("OnResponse paniced: %v", r))
			e = errors.New("OnResponse paniced")
		}
	}()

	sp := opentracing.StartSpan(fmt.Sprintf("PLUGIN: ON_RESPONSE: %s", plugin.Name), opentracing.ChildOf(span.Context()))
	defer sp.Finish()

	return p.OnResponse(ctx, plugin.Config, m, proxy, req, resp, sp)
}

func getPluginFileName(p v2.Plugin) string {
	if strings.Compare(p.Version, "") != 0 {
		return fmt.Sprintf("%s_%s",
			p.Name,
			p.Version,
		)
	}
	return p.Name
}

func getPlugin(pl v2.Plugin) (*plugin.Plugin, error) {
	path, err := getAbsPath(viper.GetString(options.FlagPluginsLocation.GetLong()))
	if err != nil {
		return nil, kanaliErrors.StatusError{Code: http.StatusInternalServerError, Err: fmt.Errorf("file path %s could not be found", viper.GetString(options.FlagPluginsLocation.GetLong()))}
	}

	plug, err := pluginPkg.Open(fmt.Sprintf("%s/%s.so",
		path,
		getPluginFileName(pl),
	))
	if err != nil {
		return nil, kanaliErrors.StatusError{
			Code: http.StatusInternalServerError,
			Err:  fmt.Errorf("could not open plugin %s: %s", pl.Name, err.Error()),
		}
	}

	symPlug, err := plug.Lookup(pluginSymbolName)
	if err != nil {
		return nil, kanaliErrors.StatusError{
			Code: http.StatusInternalServerError,
			Err:  err,
		}
	}

	var p plugin.Plugin
	p, ok := symPlug.(plugin.Plugin)
	if !ok {
		return nil, kanaliErrors.StatusError{
			Code: http.StatusInternalServerError,
			Err:  fmt.Errorf("plugin %s must implement the Plugin interface", pl.Name),
		}
	}

	return &p, nil
}

func getAbsPath(path string) (string, error) {

	p, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	if p[len(p)-1] == '/' {
		if len(p) < 2 {
			return "", nil
		}
		return p[:len(p)-2], nil
	}

	return p, nil

}

type proxyPassStep struct{}

func (step proxyPassStep) getName() string {
	return "proxy pass"
}

func (step proxyPassStep) do(ctx context.Context, proxy *v2.ApiProxy, k8sCoreClient core.Interface, m *metrics.Metrics, w http.ResponseWriter, r *http.Request, resp *http.Response, span opentracing.Span) error {

	targetRequest, err := createTargetRequest(proxy, k8sCoreClient, r)
	if err != nil {
		return err
	}

	targetClient, err := createTargetClient(proxy, k8sCoreClient, r)
	if err != nil {
		return err
	}

	targetResponse, err := preformTargetProxy(ctx, targetClient, targetRequest, m, span)
	if err != nil {
		return err
	}

	*resp = *targetResponse
	return nil

}

func createTargetRequest(proxy *v2.ApiProxy, k8sCoreClient core.Interface, originalRequest *http.Request) (*http.Request, error) {
	targetRequest := &http.Request{}
	*targetRequest = *originalRequest
	targetRequest.RequestURI = ""

	u, err := getTargetURL(proxy, k8sCoreClient, originalRequest)
	if err != nil {
		return nil, err
	}

	targetRequest.URL = u

	targetRequest.Header.Del("apikey")
	targetRequest.Header.Add("X-Forwarded-For", originalRequest.RemoteAddr)

	return targetRequest, nil
}

func createTargetClient(proxy *v2.ApiProxy, k8sCoreClient core.Interface, originalRequest *http.Request) (*http.Client, error) {
	client := &http.Client{
		Timeout: viper.GetDuration(options.FlagProxyUpstreamTimeout.GetLong()),
	}

	transport, err := configureTargetTLS(proxy, k8sCoreClient, originalRequest)
	if err != nil {
		return nil, err
	}
	if transport != nil {
		client.Transport = transport
	}

	return client, nil
}

func configureTargetTLS(proxy *v2.ApiProxy, k8sCoreClient core.Interface, originalRequest *http.Request) (*http.Transport, error) {

	secret, err := k8sCoreClient.V1().Secrets().Lister().Secrets(proxy.ObjectMeta.Namespace).Get(proxy.Spec.Target.SSL.SecretName)
	if err != nil {
		switch e := err.(type) {
		case *k8sErrors.StatusError:
			if e.ErrStatus.Reason == metav1.StatusReasonNotFound {
				return nil, nil
			}
		default:
			return nil, kanaliErrors.StatusError{Code: http.StatusInternalServerError, Err: err}
		}
	}

	if secret.Type != v1.SecretTypeTLS {
		return nil, nil
	}

	tlsConfig := &tls.Config{}
	caCertPool := x509.NewCertPool()

	// server side tls must be configured
	cert, err := x509KeyPair(secret)
	if err != nil {
		return nil, kanaliErrors.StatusError{Code: http.StatusInternalServerError, Err: err}
	}
	tlsConfig.Certificates = []tls.Certificate{*cert}

	if secret.Data["tls.ca"] != nil {
		caCertPool.AppendCertsFromPEM(secret.Data["tls.ca"])
	}

	if !viper.GetBool(options.FlagProxyTLSCommonNameValidation.GetLong()) {
		tlsConfig.InsecureSkipVerify = true
		tlsConfig.VerifyPeerCertificate = func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
			opts := x509.VerifyOptions{
				Roots: caCertPool,
			}
			cert, err := x509.ParseCertificate(rawCerts[0])
			if err != nil {
				return err
			}
			_, err = cert.Verify(opts)
			return err
		}
	}

	tlsConfig.RootCAs = caCertPool
	tlsConfig.BuildNameToCertificate()
	return &http.Transport{TLSClientConfig: tlsConfig}, nil

}

// x509KeyPair creates a tls.Certificate from the tls data in
// a Kubernetes secret of type kubernetes.io/tls
func x509KeyPair(s *v1.Secret) (*tls.Certificate, error) {
	pair, err := tls.X509KeyPair(s.Data["tls.crt"], s.Data["tls.key"])
	if err != nil {
		return nil, err
	}
	return &pair, err
}

func preformTargetProxy(ctx context.Context, client httpClient, request *http.Request, m *metrics.Metrics, span opentracing.Span) (*http.Response, error) {
	logger := logging.WithContext(ctx)

	if err := span.Tracer().Inject(
		span.Context(),
		opentracing.TextMap,
		opentracing.HTTPHeadersCarrier(request.Header),
	); err != nil {
		logger.Error(err.Error())
	}

	sp := opentracing.StartSpan(fmt.Sprintf("%s %s",
		request.Method,
		request.URL.EscapedPath(),
	), opentracing.ChildOf(span.Context()))
	defer sp.Finish()

	hydrateSpanFromRequest(request, sp)

	t0 := time.Now()
	resp, err := client.Do(request)
	if err != nil {
		return nil, kanaliErrors.StatusError{Code: http.StatusInternalServerError, Err: err}
	}

	m.Add(
		metrics.Metric{Name: "total_target_time", Value: int(time.Now().Sub(t0) / time.Millisecond), Index: false},
	)

	hydrateSpanFromResponse(resp, sp)

	return resp, nil
}

func getTargetURL(proxy *v2.ApiProxy, k8sCoreClient core.Interface, originalRequest *http.Request) (*url.URL, error) {

	scheme := "http"

	if len(proxy.Spec.Target.SSL.SecretName) > 0 {
		scheme += "s"
	}

	services, err := k8sCoreClient.V1().Services().Lister().Services(proxy.ObjectMeta.Namespace).List(labels.SelectorFromSet(getServiceLabelSet(proxy, originalRequest.Header)))
	if err != nil || len(services) == 0 {
		return nil, kanaliErrors.StatusError{Code: http.StatusNotFound, Err: errors.New("no matching services")}
	}

	uri := fmt.Sprintf("%s.%s.svc.cluster.local",
		services[0].ObjectMeta.Name,
		proxy.ObjectMeta.Namespace,
	)

	if viper.GetBool(options.FlagProxyEnableClusterIP.GetLong()) {
		uri = services[0].Spec.ClusterIP
	}

	return &url.URL{
		Scheme: scheme,
		Host: fmt.Sprintf("%s:%d",
			uri,
			proxy.Spec.Target.Backend.Service.Port,
		),
		Path:       utils.ComputeTargetPath(proxy.Spec.Source.Path, proxy.Spec.Target.Path, originalRequest.URL.EscapedPath()),
		RawPath:    utils.ComputeTargetPath(proxy.Spec.Source.Path, proxy.Spec.Target.Path, originalRequest.URL.EscapedPath()),
		ForceQuery: originalRequest.URL.ForceQuery,
		RawQuery:   originalRequest.URL.RawQuery,
		Fragment:   originalRequest.URL.Fragment,
	}, nil

}

func getServiceLabelSet(proxy *v2.ApiProxy, requestHeaders http.Header) labels.Set {
	ls := map[string]string{}
	for _, label := range proxy.Spec.Target.Backend.Service.Labels {
		if len(label.Header) > 0 {
			ls[label.Name] = requestHeaders.Get(label.Header)
		} else {
			ls[label.Name] = label.Value
		}
	}
	return ls
}

// ValidateProxyStep is factory that defines a step responsible for
// validating that an incoming request matches a proxy that Kanali
// has stored in memory
type validateProxyStep struct{}

// GetName retruns the name of the ValidateProxyStep step
func (step validateProxyStep) getName() string {
	return "Validate Proxy"
}

// Do executes the logic of the ValidateProxyStep step
func (step validateProxyStep) do(ctx context.Context, proxy *v2.ApiProxy, k8sCoreClient core.Interface, m *metrics.Metrics, w http.ResponseWriter, r *http.Request, resp *http.Response, trace opentracing.Span) error {

	untypedProxy, err := store.ApiProxyStore.Get(utils.ComputeURLPath(r.URL))
	if err != nil || untypedProxy == nil {
		if err != nil {
			logging.WithContext(ctx).Error("error retrieving proxy", zap.String("msg", err.Error()))
		}

		trace.SetTag(tags.KanaliProxyName, "unknown")
		trace.SetTag(tags.KanaliProxyNamespace, "unknown")

		m.Add(
			metrics.Metric{Name: "proxy_name", Value: "unknown", Index: true},
			metrics.Metric{Name: "proxy_namespace", Value: "unknown", Index: true},
		)

		return kanaliErrors.StatusError{Code: http.StatusNotFound, Err: errors.New("proxy not found")}
	}

	typedProxy, _ := untypedProxy.(v2.ApiProxy)

	*proxy = typedProxy

	trace.SetTag(tags.KanaliProxyName, proxy.ObjectMeta.Name)
	trace.SetTag(tags.KanaliProxyNamespace, proxy.ObjectMeta.Namespace)

	m.Add(
		metrics.Metric{Name: "proxy_name", Value: proxy.ObjectMeta.Name, Index: true},
		metrics.Metric{Name: "proxy_namespace", Value: proxy.ObjectMeta.Namespace, Index: true},
	)

	return nil

}

// WriteResponseStep is factory that defines a step responsible for writing
// an HTTP response
type writeResponseStep struct{}

// GetName retruns the name of the WriteResponseStep step
func (step writeResponseStep) getName() string {
	return "Write Response"
}

// Do executes the logic of the WriteResponseStep step
func (step writeResponseStep) do(ctx context.Context, proxy *v2.ApiProxy, k8sCoreClient core.Interface, m *metrics.Metrics, w http.ResponseWriter, r *http.Request, resp *http.Response, span opentracing.Span) error {

	for k, v := range resp.Header {
		for _, value := range v {
			w.Header().Set(k, value)
		}
	}

	m.Add(metrics.Metric{Name: "http_response_code", Value: strconv.Itoa(resp.StatusCode), Index: true})

	hydrateSpanFromResponse(resp, span)

	w.WriteHeader(resp.StatusCode)

	if _, err := io.Copy(w, resp.Body); err != nil {
		return err
	}

	return nil
}
