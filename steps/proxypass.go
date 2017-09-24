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

package steps

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/northwesternmutual/kanali/config"
	"github.com/northwesternmutual/kanali/metrics"
	"github.com/northwesternmutual/kanali/spec"
	"github.com/northwesternmutual/kanali/tracer"
	"github.com/northwesternmutual/kanali/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/spf13/viper"
	"k8s.io/kubernetes/pkg/api"
)

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// ProxyPassStep is factory that defines a step responsible for configuring
// and performing a proxy to a dynamic upstream service
type ProxyPassStep struct{}

// GetName retruns the name of the ProxyPassStep step
func (step ProxyPassStep) GetName() string {
	return "Proxy Pass"
}

// Do executes the logic of the ProxyPassStep step
func (step ProxyPassStep) Do(ctx context.Context, proxy *spec.APIProxy, m *metrics.Metrics, w http.ResponseWriter, r *http.Request, resp *http.Response, span opentracing.Span) error {

	targetRequest, err := createTargetRequest(proxy, r)
	if err != nil {
		return err
	}

	targetClient, err := createTargetClient(proxy, r)
	if err != nil {
		return err
	}

	targetResponse, err := preformTargetProxy(targetClient, targetRequest, m, span)
	if err != nil {
		return err
	}

	*resp = *targetResponse
	return nil

}

func createTargetRequest(proxy *spec.APIProxy, originalRequest *http.Request) (*http.Request, error) {
	targetRequest := &http.Request{}
	*targetRequest = *originalRequest
	targetRequest.RequestURI = ""

	u, err := getTargetHost(proxy, originalRequest)
	if err != nil {
		return nil, err
	}

	u.Path = utils.ComputeTargetPath(proxy.Spec.Path, proxy.Spec.Target, originalRequest.URL.Path)
	u.RawPath = utils.ComputeTargetPath(proxy.Spec.Path, proxy.Spec.Target, originalRequest.URL.EscapedPath())
	u.ForceQuery = originalRequest.URL.ForceQuery
	u.RawQuery = originalRequest.URL.RawQuery
	u.Fragment = originalRequest.URL.Fragment

	targetRequest.URL = u

	targetRequest.Header.Del("apikey")
	targetRequest.Header.Add("X-Forwarded-For", originalRequest.RemoteAddr)

	return targetRequest, nil
}

func createTargetClient(proxy *spec.APIProxy, originalRequest *http.Request) (*http.Client, error) {
	client := &http.Client{
		Timeout: viper.GetDuration(config.FlagProxyUpstreamTimeout.GetLong()),
	}

	transport, err := configureTargetTLS(proxy, originalRequest)
	if err != nil {
		return nil, err
	}

	client.Transport = transport

	return client, nil
}

func configureTargetTLS(proxy *spec.APIProxy, originalRequest *http.Request) (*http.Transport, error) {

	// get secret for this request - if any
	untypedSecret, err := spec.SecretStore.Get(proxy.GetSSLCertificates(originalRequest.Host).SecretName, proxy.ObjectMeta.Namespace)
	if err != nil {
		return nil, utils.StatusError{Code: http.StatusInternalServerError, Err: err}
	}

	tlsConfig := &tls.Config{}
	caCertPool := x509.NewCertPool()

	// ssl is not configured for this request
	if untypedSecret == nil {

		// if upstream option is being used, if the scheme
		// is https we need to add the root ca bundle
		if strings.Compare(originalRequest.URL.Scheme, "https") != 0 {
			logrus.Debug("TLS not configured for this proxy")
			return nil, nil
		}

	} else {

		secret, ok := untypedSecret.(api.Secret)
		if !ok {
			return nil, utils.StatusError{Code: http.StatusInternalServerError, Err: errors.New("the secret store is corrupted")}
		}

		// server side tls must be configured
		cert, err := spec.X509KeyPair(secret)
		if err != nil {
			return nil, utils.StatusError{Code: http.StatusInternalServerError, Err: err}
		}
		tlsConfig.Certificates = []tls.Certificate{*cert}

		if secret.Data["tls.ca"] != nil {
			caCertPool.AppendCertsFromPEM(secret.Data["tls.ca"])
		}

		if !viper.GetBool(config.FlagProxyTLSCommonNameValidation.GetLong()) {
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

	}

	tlsConfig.RootCAs = caCertPool
	tlsConfig.BuildNameToCertificate()
	return &http.Transport{TLSClientConfig: tlsConfig}, nil

}

func preformTargetProxy(client httpClient, request *http.Request, m *metrics.Metrics, span opentracing.Span) (*http.Response, error) {
	if err := span.Tracer().Inject(
		span.Context(),
		opentracing.TextMap,
		opentracing.HTTPHeadersCarrier(request.Header),
	); err != nil {
		logrus.Error("error injecting headers")
	}

	sp := opentracing.StartSpan(fmt.Sprintf("%s %s",
		request.Method,
		request.URL.EscapedPath(),
	), opentracing.ChildOf(span.Context()))
	defer sp.Finish()

	tracer.HydrateSpanFromRequest(request, sp)

	t0 := time.Now()
	resp, err := client.Do(request)
	if err != nil {
		return nil, utils.StatusError{Code: http.StatusInternalServerError, Err: err}
	}

	m.Add(
		metrics.Metric{Name: "total_target_time", Value: int(time.Now().Sub(t0) / time.Millisecond), Index: false},
	)

	tracer.HydrateSpanFromResponse(resp, sp)

	return resp, nil
}

func getTargetHost(proxy *spec.APIProxy, originalRequest *http.Request) (*url.URL, error) {

	scheme := "http"

	if *proxy.GetSSLCertificates(originalRequest.Host) != (spec.SSL{}) {
		scheme = "https"
	}

	untypedSvc, err := spec.ServiceStore.Get(proxy.Spec.Service, originalRequest.Header)
	if err != nil || untypedSvc == nil {
		logrus.Debug("service was non of type spec.Service")
		return nil, utils.StatusError{Code: http.StatusNotFound, Err: errors.New("no matching services")}
	}

	svc, _ := untypedSvc.(spec.Service)

	uri := fmt.Sprintf("%s.%s.svc.cluster.local",
		svc.Name,
		proxy.ObjectMeta.Namespace,
	)

	if viper.GetBool(config.FlagProxyEnableClusterIP.GetLong()) {
		uri = svc.ClusterIP
	}

	return &url.URL{
		Scheme: scheme,
		Host: fmt.Sprintf("%s:%d",
			uri,
			proxy.Spec.Service.Port,
		),
	}, nil

}
