package steps

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/northwesternmutual/kanali/cmd/kanali/app/options"
	"github.com/northwesternmutual/kanali/pkg/apis/kanali.io/v2"
	kanaliErrors "github.com/northwesternmutual/kanali/pkg/errors"
	"github.com/northwesternmutual/kanali/pkg/logging"
	"github.com/northwesternmutual/kanali/pkg/metrics"
	"github.com/northwesternmutual/kanali/pkg/tags"
	"github.com/northwesternmutual/kanali/pkg/tracer"
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

// httpClient allows for mocking an http client to assist in testing
type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type proxyPassStep struct{}

func NewProxyPassStep() Step {
	return proxyPassStep{}
}

func (step proxyPassStep) GetName() string {
	return "proxy pass"
}

func (step proxyPassStep) Do(ctx context.Context, proxy *v2.ApiProxy, k8sCoreClient core.Interface, m *metrics.Metrics, w http.ResponseWriter, r *http.Request, resp *http.Response, span opentracing.Span) error {

	targetRequest, err := createTargetRequest(ctx, proxy, k8sCoreClient, r)
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

func createTargetRequest(ctx context.Context, proxy *v2.ApiProxy, k8sCoreClient core.Interface, originalRequest *http.Request) (*http.Request, error) {
	targetRequest := &http.Request{}
	*targetRequest = *originalRequest
	targetRequest.RequestURI = ""

	u, err := getTargetURL(ctx, proxy, k8sCoreClient, originalRequest)
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

	// if len(proxy.Spec.Target.Backend.Endpoint) > 0 { // Endpoint backend is being used
	//   tlsConfig := &tls.Config{}
	//   caCertPool := x509.NewCertPool()
	//   tlsConfig.RootCAs = caCertPool
	//   tlsConfig.InsecureSkipVerify = true
	// 	tlsConfig.BuildNameToCertificate()
	// 	return &http.Transport{TLSClientConfig: tlsConfig}, nil
	// }

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

	tracer.HydrateSpanFromRequest(request, sp)

	logger.With(
		zap.String(tags.HTTPRequestURLScheme, request.URL.Scheme),
		zap.String(tags.HTTPRequestURLHost, request.URL.Host),
		zap.String(tags.HTTPRequestURLPath, request.URL.Path),
	).Info("upstream request")

	t0 := time.Now()
	resp, err := client.Do(request)
	t1 := time.Now()
	if err != nil {
		return nil, kanaliErrors.StatusError{Code: http.StatusInternalServerError, Err: err}
	}

	m.Add(
		metrics.Metric{Name: "total_target_time", Value: int(t1.Sub(t0) / time.Millisecond), Index: false},
	)

	tracer.HydrateSpanFromResponse(resp, sp)

	return resp, nil
}

func getTargetURL(ctx context.Context, proxy *v2.ApiProxy, k8sCoreClient core.Interface, originalRequest *http.Request) (*url.URL, error) {

	scheme, host := "", ""
	logger := logging.WithContext(ctx)

	if len(proxy.Spec.Target.Backend.Endpoint) > 0 { // Endpoint backend is configured

		url, err := url.Parse(proxy.Spec.Target.Backend.Endpoint)
		if err != nil {
			logger.Error(err.Error())
			return nil, kanaliErrors.StatusError{Code: http.StatusInternalServerError, Err: errors.New("error parsing upstream url")}
		}
		scheme = url.Scheme
		host = url.Host

	} else { // Service backend is configured
		if len(proxy.Spec.Target.SSL.SecretName) > 0 {
			scheme = "https"
		} else {
			scheme = "http"
		}

		services, err := k8sCoreClient.V1().Services().Lister().Services(proxy.ObjectMeta.Namespace).List(labels.SelectorFromSet(getServiceLabelSet(proxy, originalRequest.Header)))
		if err != nil {
			logger.Error(err.Error())
		}
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

		host = fmt.Sprintf("%s:%d",
			uri,
			proxy.Spec.Target.Backend.Service.Port,
		)
	}

	return &url.URL{
		Scheme:     scheme,
		Host:       host,
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
