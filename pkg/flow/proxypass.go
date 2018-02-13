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
	"strconv"
	"strings"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"k8s.io/api/core/v1"
	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"

	"github.com/northwesternmutual/kanali/cmd/kanali/app/options"
	"github.com/northwesternmutual/kanali/pkg/apis/kanali.io/v2"
	"github.com/northwesternmutual/kanali/pkg/errors"
	"github.com/northwesternmutual/kanali/pkg/log"
	coreV1 "github.com/northwesternmutual/kanali/pkg/store/core/v1"
	store "github.com/northwesternmutual/kanali/pkg/store/kanali/v2"
	"github.com/northwesternmutual/kanali/pkg/tags"
	"github.com/northwesternmutual/kanali/pkg/tracer"
	"github.com/northwesternmutual/kanali/pkg/utils"
)

// httpClient allows for mocking an http client to assist in testing
type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type proxyPassStep struct {
	logger               *zap.Logger
	span                 opentracing.Span
	proxy                *v2.ApiProxy
	originalReq          *http.Request
	upstreamReq          *http.Request
	upstreamRoundTripper http.RoundTripper
	upstreamResp         *http.Response
	originalRespWriter   http.ResponseWriter
	err                  error
}

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

func ProxyPassStep() Step {
	return proxyPassStep{}
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

	// TODO: set timeout

	transport := http.DefaultTransport.(*http.Transport)
	tlsConfig, err := step.configureTLS()
	if err != nil {
		step.err = err
		return step
	}

	transport.TLSClientConfig = tlsConfig
	step.upstreamRoundTripper = transport

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

	copyBuffer(step.originalReq.Context(), step.originalRespWriter, step.upstreamResp.Body, []byte{})
	step.upstreamResp.Body.Close() // close now, instead of defer, to populate res.Trailer

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

	secret, err := coreV1.Interface().Secrets().Lister().Secrets(step.proxy.GetNamespace()).Get(step.proxy.Spec.Target.SSL.SecretName)
	if err != nil {
		switch e := err.(type) {
		case *k8sErrors.StatusError:
			if e.ErrStatus.Reason == metav1.StatusReasonNotFound {
				return nil, nil
			}
		default:
			logger.Error(err.Error())
			return nil, errors.ErrorKubernetesSecretError
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
		logger.Error(err.Error())
		return nil, errors.ErrorCreateKeyPair
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
	return tlsConfig, nil
}

func (step proxyPassStep) serviceDetails() (string, string, error) {
	logger := log.WithContext(step.originalReq.Context())
	var scheme string

	if len(step.proxy.Spec.Target.SSL.SecretName) > 0 {
		scheme = "https"
	} else {
		scheme = "http"
	}

	services, err := coreV1.Interface().Services().Lister().Services(step.proxy.GetNamespace()).List(labels.SelectorFromSet(getServiceLabelSet(step.proxy, step.originalReq.Header)))
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
		step.upstreamReq.URL.Scheme, step.upstreamReq.URL.Host = step.proxy.Spec.Target.Backend.Endpoint.Scheme, step.proxy.Spec.Target.Backend.Endpoint.Host
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

func getServiceLabelSet(proxy *v2.ApiProxy, reqHeaders http.Header) labels.Set {
	ls := map[string]string{}
	for _, label := range proxy.Spec.Target.Backend.Service.Labels {
		if len(label.Header) > 0 {
			ls[label.Name] = reqHeaders.Get(label.Header)
		} else {
			ls[label.Name] = label.Value
		}
	}
	return ls
}
