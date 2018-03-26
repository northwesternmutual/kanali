// Copyright (c) 2017 Northwestern Mutual.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package flow

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"k8s.io/api/core/v1"
	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	corev1 "k8s.io/client-go/informers/core/v1"

	"github.com/northwesternmutual/kanali/cmd/kanali/app/options"
	"github.com/northwesternmutual/kanali/pkg/apis/kanali.io/v2"
	"github.com/northwesternmutual/kanali/pkg/errors"
	"github.com/northwesternmutual/kanali/pkg/log"
	store "github.com/northwesternmutual/kanali/pkg/store/kanali/v2"
	"github.com/northwesternmutual/kanali/pkg/tags"
	"github.com/northwesternmutual/kanali/pkg/tracer"
	"github.com/northwesternmutual/kanali/pkg/utils"
	tlsutils "github.com/northwesternmutual/kanali/pkg/utils/tls"
)

// httpClient allows for mocking an http client to assist in testing
type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type proxyPassStep struct {
	logger               *zap.Logger
	v1Interface          corev1.Interface
	span                 opentracing.Span
	proxy                *v2.ApiProxy
	originalReq          *http.Request
	upstreamReq          *http.Request
	upstreamRoundTripper http.RoundTripper
	upstreamResp         *http.Response
	originalRespWriter   http.ResponseWriter
	err                  error
}

var (
	defaultTLSCertKey = "tls.crt"
	defaultTLSKeyKey  = "tls.key"
	defaultTLSCAKey   = "tls.ca"
)

// Hop-by-hop headers. These are removed when sent to the backend.
// http://www.w3.org/Protocols/rfc2616/rfc2616-sec13.html
var hopHeaders = []string{
	"Connection",
	"Proxy-Connection", // non-standard but still sent by libcurl and rejected by e.g. google
	"Keep-Alive",
	"Proxy-Authenticate",
	"Proxy-Authorization",
	"Te",      // canonicalized version of "TE"
	"Trailer", // not Trailers per URL above; http://www.rfc-editor.org/errata_search.php?eid=4522
	"Transfer-Encoding",
	"Upgrade",
}

func ProxyPassStep(i corev1.Interface) Step {
	return proxyPassStep{
		v1Interface: i,
	}
}

func (step proxyPassStep) Name() string {
	return "Proxy Pass"
}

func (step proxyPassStep) Do(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	logger := log.WithContext(r.Context())

	proxy := store.ApiProxyStore().Get(utils.ComputeURLPath(r.URL))
	if proxy == nil {
		logger.Warn(errors.ErrorProxyNotFound.Message)
		return errors.ErrorProxyNotFound
	}

	step.logger = logger
	step.span = opentracing.SpanFromContext(ctx)
	step.originalReq = r
	step.originalRespWriter = w
	step.proxy = proxy

	results := step.configureRequest().configureTransport().preformProxyPass().writeResponse()
	if results.err != nil {
		logger.Error(results.err.Error())
		return results.err
	}

	return next()
}

func (step proxyPassStep) configureRequest() proxyPassStep {
	if step.err != nil {
		return step
	}

	ctx := step.originalReq.Context()
	// preforms a shallow copy of http.Request with a
	// deep copy of http.Request.URL
	step.upstreamReq = step.originalReq.WithContext(ctx)

	// Remove hop-by-hop headers listed in the "Connection" header.
	// See RFC 2616, section 14.10.
	if c := step.upstreamReq.Header.Get("Connection"); c != "" {
		for _, f := range strings.Split(c, ",") {
			if f = strings.TrimSpace(f); f != "" {
				step.upstreamReq.Header.Del(f)
			}
		}
	}

	// Remove hop-by-hop headers to the backend. Especially
	// important is "Connection" because we want a persistent
	// connection, regardless of what the client sent to us.
	for _, h := range hopHeaders {
		if step.upstreamReq.Header.Get(h) != "" {
			step.upstreamReq.Header.Del(h)
		}
	}

	step.upstreamReq.Header = utils.CloneHTTPHeader(step.originalReq.Header)
	step.upstreamReq.URL.Host = step.originalReq.Host
	if _, ok := step.upstreamReq.Header["User-Agent"]; !ok {
		// explicitly disable User-Agent so it's not set to default value
		step.upstreamReq.Header.Set("User-Agent", "")
	}
	step.upstreamReq.Close = false

	if clientIP, _, err := net.SplitHostPort(step.originalReq.RemoteAddr); err == nil {
		// If we aren't the first proxy retain prior
		// X-Forwarded-For information as a comma+space
		// separated list and fold multiple headers into one.
		if prior, ok := step.upstreamReq.Header["X-Forwarded-For"]; ok {
			clientIP = strings.Join(prior, ", ") + ", " + clientIP
		}
		step.upstreamReq.Header.Set("X-Forwarded-For", clientIP)
	}

	// TODO: should not be hardcoded!
	step.upstreamReq.Header.Del("apikey")

	if err := step.setUpstreamURL(); err != nil {
		step.err = err
	}
	return step
}

func (step proxyPassStep) configureTransport() proxyPassStep {
	if step.err != nil {
		return step
	}

	//transport := http.DefaultTransport.(*http.Transport)
	tlsConfig, err := step.configureTLS()
	if err != nil {
		step.err = err
		return step
	}

	//transport.TLSClientConfig = tlsConfig
	step.upstreamRoundTripper = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   100 * time.Second,
		ExpectContinueTimeout: 10 * time.Second,
		TLSClientConfig:       tlsConfig,
	}

	return step
}

func (step proxyPassStep) preformProxyPass() proxyPassStep {
	if step.err != nil {
		return step
	}

	logger := log.WithContext(step.originalReq.Context())
	var upstreamSpan opentracing.Span

	if step.span != nil {
		if err := step.span.Tracer().Inject(
			step.span.Context(),
			opentracing.TextMap,
			opentracing.HTTPHeadersCarrier(step.upstreamReq.Header),
		); err != nil {
			logger.Error(err.Error())
		}

		upstreamSpan = opentracing.StartSpan(fmt.Sprintf("%s %s",
			step.upstreamReq.Method,
			step.upstreamReq.URL.EscapedPath(),
		), opentracing.ChildOf(step.span.Context()))
		defer upstreamSpan.Finish()

		tracer.HydrateSpanFromRequest(step.upstreamReq, upstreamSpan)
	}

	logger.Info("upstream request",
		zap.String(tags.HTTPRequestURLScheme, step.upstreamReq.URL.Scheme),
		zap.String(tags.HTTPRequestURLHost, step.upstreamReq.URL.Host),
		zap.String(tags.HTTPRequestURLPath, step.upstreamReq.URL.Path),
	)

	res, err := step.upstreamRoundTripper.RoundTrip(step.upstreamReq)
	if err != nil {
		logger.Error(err.Error())
		step.err = errors.ErrorBadGateway
		return step
	}

	if step.span != nil {
		tracer.HydrateSpanFromResponse(res, upstreamSpan)
	}

	// Remove hop-by-hop headers listed in the
	// "Connection" header of the response.
	if c := res.Header.Get("Connection"); c != "" {
		for _, f := range strings.Split(c, ",") {
			if f = strings.TrimSpace(f); f != "" {
				res.Header.Del(f)
			}
		}
	}

	for _, h := range hopHeaders {
		res.Header.Del(h)
	}

	step.upstreamResp = res

	return step
}

func (step proxyPassStep) writeResponse() proxyPassStep {
	if step.err != nil {
		return step
	}

	logger := log.WithContext(step.originalReq.Context())

	copyHeader(step.originalRespWriter.Header(), step.upstreamResp.Header)

	// The "Trailer" header isn't included in the Transport's response,
	// at least for *http.Transport. Build it up from Trailer.
	announcedTrailers := len(step.upstreamResp.Trailer)
	if announcedTrailers > 0 {
		trailerKeys := make([]string, 0, len(step.upstreamResp.Trailer))
		for k := range step.upstreamResp.Trailer {
			trailerKeys = append(trailerKeys, k)
		}
		step.originalRespWriter.Header().Add("Trailer", strings.Join(trailerKeys, ", "))
	}

	step.originalRespWriter.WriteHeader(step.upstreamResp.StatusCode)
	if len(step.upstreamResp.Trailer) > 0 {
		// Force chunking if we saw a response trailer.
		// This prevents net/http from calculating the length for short
		// bodies and adding a Content-Length.
		if fl, ok := step.originalRespWriter.(http.Flusher); ok {
			fl.Flush()
		}
	}

	if n, err := copyBuffer(step.originalReq.Context(), step.originalRespWriter, step.upstreamResp.Body, []byte{}); err != nil {
		logger.Warn(fmt.Sprintf("error writing to response - %d bytes written: %s", n, err))
	}
	if err := step.upstreamResp.Body.Close(); err != nil { // close now, instead of defer, to populate res.Trailer
		logger.Warn(fmt.Sprintf("error closing response body: %s", err))
	}

	if len(step.upstreamResp.Trailer) == announcedTrailers {
		copyHeader(step.originalRespWriter.Header(), step.upstreamResp.Trailer)
		return step
	}

	for k, vv := range step.upstreamResp.Trailer {
		k = http.TrailerPrefix + k
		for _, v := range vv {
			step.originalRespWriter.Header().Add(k, v)
		}
	}

	return step
}

func (step proxyPassStep) configureTLS() (*tls.Config, error) {
	logger := log.WithContext(step.originalReq.Context())

	// determine whether tls is needed at all
	if strings.ToLower(step.upstreamReq.URL.Scheme) != "https" {
		logger.Debug(fmt.Sprintf("tls is not needed for scheme %s", step.upstreamReq.URL.Scheme))
		return nil, nil
	}

	config := new(tls.Config)

	if step.proxy.Spec.Target.SSL != nil && len(step.proxy.Spec.Target.SSL.SecretName) > 0 {
		config.RootCAs = x509.NewCertPool()
		secret, err := step.v1Interface.Secrets().Lister().Secrets(step.proxy.GetNamespace()).Get(step.proxy.Spec.Target.SSL.SecretName)
		if err != nil {
			err := fmt.Errorf("secret %s not found in %s namesapce", step.proxy.Spec.Target.SSL.SecretName, step.proxy.GetNamespace())
			logger.Error(err.Error())
			return nil, err
		}

		if !metav1.HasAnnotation(secret.ObjectMeta, "kanali.io/enabled") || secret.ObjectMeta.GetAnnotations()["kanali.io/enabled"] != "true" {
			err := fmt.Errorf("secret %s in %s namespaces exists - however, due to the annotations, kanali doesn't have permission to use this secret", step.proxy.Spec.Target.SSL.SecretName, step.proxy.GetNamespace())
			logger.Info(err.Error())
			return nil, err
		}

		cert, key := getCertKey(secret)
		ca := getCA(secret)

		if cert == nil && key == nil && ca == nil {
			return nil, fmt.Errorf("secret does not contain any valid data")
		}

		if cert != nil && key != nil && len(cert) > 0 && len(key) > 0 {
			pair, err := tls.X509KeyPair(cert, key)
			if err != nil {
				logger.Error(err.Error())
				return nil, errors.ErrorCreateKeyPair
			}
			config.Certificates = []tls.Certificate{pair}
		}

		if ca != nil {
			if ok := config.RootCAs.AppendCertsFromPEM(ca); !ok {
				return nil, fmt.Errorf("could not append certificate to pool")
			}
		}
	} else {
		pool, err := tlsutils.GetSystemCertPool()
		if err != nil {
			return nil, err
		}
		config.RootCAs = pool
	}

	// check if common name or sans validation should be performed
	if !viper.GetBool(options.FlagProxyTLSCommonNameValidation.GetLong()) {
		config.InsecureSkipVerify = true
		config.VerifyPeerCertificate = tlsutils.VerifyPeerCertificate(config.RootCAs)
	}

	config.BuildNameToCertificate()
	return config, nil
}

func getCertKey(secret *v1.Secret) (cert, key []byte) {
	if metav1.HasAnnotation(secret.ObjectMeta, "kanali.io/cert") {
		cert = secret.Data[secret.GetAnnotations()["kanali.io/cert"]]
	} else {
		cert = secret.Data[defaultTLSCertKey]
	}

	if metav1.HasAnnotation(secret.ObjectMeta, "kanali.io/key") {
		key = secret.Data[secret.GetAnnotations()["kanali.io/key"]]
	} else {
		key = secret.Data[defaultTLSKeyKey]
	}

	return
}

func getCA(secret *v1.Secret) []byte {
	if metav1.HasAnnotation(secret.ObjectMeta, "kanali.io/ca") {
		return secret.Data[secret.GetAnnotations()["kanali.io/ca"]]
	}
	return secret.Data[defaultTLSCAKey]
}

func (step proxyPassStep) serviceDetails() (string, string, error) {
	logger := log.WithContext(step.originalReq.Context())
	var scheme string

	if step.proxy.Spec.Target.SSL != nil && len(step.proxy.Spec.Target.SSL.SecretName) > 0 {
		scheme = "https"
	} else {
		scheme = "http"
	}

	services, err := step.v1Interface.Services().Lister().Services(step.proxy.GetNamespace()).List(labels.SelectorFromSet(
		getServiceLabelSet(step.proxy, step.originalReq.Header, viper.GetStringMapString(options.FlagProxyDefaultHeaderValues.GetLong())),
	))
	if err != nil || len(services) < 1 {
		switch e := err.(type) {
		case *k8sErrors.StatusError:
			if e.ErrStatus.Reason == metav1.StatusReasonNotFound {
				return "", "", errors.ErrorNoMatchingServices
			}
		default:
			logger.Error(err.Error())
			return "", "", errors.ErrorKubernetesServiceError
		}
	}

	if len(services) > 1 {
		step.originalRespWriter.Header().Add("x-kanali-matched-services", strconv.Itoa(len(services)))
		logger.Debug(fmt.Sprintf("there were %d matching services", len(services)))
	}

	uri := fmt.Sprintf("%s.%s.svc.cluster.local",
		services[0].GetName(),
		step.proxy.GetNamespace(),
	)

	if viper.GetBool(options.FlagProxyEnableClusterIP.GetLong()) {
		uri = services[0].Spec.ClusterIP
	}

	return scheme, fmt.Sprintf("%s:%d",
		uri,
		step.proxy.Spec.Target.Backend.Service.Port,
	), nil
}

func (step proxyPassStep) setUpstreamURL() error {
	var err error

	if step.proxy.Spec.Target.Backend.Endpoint != nil { // Endpoint backend is configured
		u, err := url.Parse(*step.proxy.Spec.Target.Backend.Endpoint)
		if err != nil {
			return err
		} else {
			step.upstreamReq.URL.Scheme = u.Scheme
			step.upstreamReq.URL.Host = u.Host
		}
	} else {
		step.upstreamReq.URL.Scheme, step.upstreamReq.URL.Host, err = step.serviceDetails()
	}

	if err != nil {
		return err
	}

	step.upstreamReq.URL.Path = utils.ComputeTargetPath(step.proxy.Spec.Source.Path, step.proxy.Spec.Target.Path, step.originalReq.URL.EscapedPath())
	if step.originalReq.URL.RawQuery == "" || step.upstreamReq.URL.RawQuery == "" {
		step.upstreamReq.URL.RawQuery = step.originalReq.URL.RawQuery + step.upstreamReq.URL.RawQuery
	} else {
		step.upstreamReq.URL.RawQuery = step.originalReq.URL.RawQuery + "&" + step.upstreamReq.URL.RawQuery
	}

	return nil
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

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func copyBuffer(ctx context.Context, dst io.Writer, src io.Reader, buf []byte) (int64, error) {
	logger := log.WithContext(ctx)
	if len(buf) == 0 {
		buf = make([]byte, 32*1024)
	}
	var written int64
	for {
		nr, rerr := src.Read(buf)
		if rerr != nil && rerr != io.EOF && rerr != context.Canceled {
			logger.Warn(fmt.Sprintf("read error during body copy: %v", rerr))
		}
		if nr > 0 {
			nw, werr := dst.Write(buf[:nr])
			if nw > 0 {
				written += int64(nw)
			}
			if werr != nil {
				return written, werr
			}
			if nr != nw {
				return written, io.ErrShortWrite
			}
		}
		if rerr != nil {
			return written, rerr
		}
	}
}

func getServiceLabelSet(p *v2.ApiProxy, h http.Header, defaults map[string]string) labels.Set {
	ls := make(map[string]string, len(p.Spec.Target.Backend.Service.Labels))
	for _, label := range p.Spec.Target.Backend.Service.Labels {
		if len(label.Header) > 0 {
			val := h.Get(label.Header)
			if val == "" {
				val = defaults[label.Header]
			}
			ls[label.Name] = val
		} else {
			ls[label.Name] = label.Value
		}
	}
	return ls
}
