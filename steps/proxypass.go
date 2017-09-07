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

	"github.com/Sirupsen/logrus"
	"github.com/northwesternmutual/kanali/controller"
	"github.com/northwesternmutual/kanali/metrics"
	"github.com/northwesternmutual/kanali/spec"
	"github.com/northwesternmutual/kanali/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/spf13/viper"
	"k8s.io/kubernetes/pkg/api"
)

type proxy struct {
	Source *http.Request
	Target spec.APIProxy
}

type upstream struct {
	Request  *http.Request
	Response *http.Response
	Client   *http.Client
	Error    utils.StatusError
}

// ProxyPassStep is factory that defines a step responsible for configuring
// and performing a proxy to a dynamic upstream service
type ProxyPassStep struct{}

// GetName retruns the name of the ProxyPassStep step
func (step ProxyPassStep) GetName() string {
	return "Proxy Pass"
}

// Do executes the logic of the ProxyPassStep step
func (step ProxyPassStep) Do(ctx context.Context, m *metrics.Metrics, c *controller.Controller, w http.ResponseWriter, r *http.Request, resp *http.Response, trace opentracing.Span) error {

	untypedProxy, err := spec.ProxyStore.Get(r.URL.Path)
	if err != nil || untypedProxy == nil {
		if err != nil {
			logrus.Error(err.Error())
		}
		return utils.StatusError{Code: http.StatusNotFound, Err: errors.New("proxy not found")}
	}

	typedProxy, ok := untypedProxy.(spec.APIProxy)
	if !ok {
		return utils.StatusError{Code: http.StatusNotFound, Err: errors.New("proxy not found")}
	}

	p := &proxy{
		Source: r,
		Target: typedProxy, // shouldn't be nil (unless the proxy is removed within the microseconds it takes to get to this code)
	}

	up := create(p).setUpstreamURL(p).configureTLS(p).setUpstreamHeaders(p).performProxy(trace)

	if up.Error != (utils.StatusError{}) {
		logrus.Errorf("error performing proxypass: %s", up.Error)
		return up.Error
	}

	*resp = *(up.Response)

	return nil

}

func create(p *proxy) *upstream {
	new := &http.Request{}
	*new = *(p.Source)

	up := &upstream{
		Request: new,
		Client: &http.Client{
			Timeout: viper.GetDuration("upstream-timeout"),
		},
		Error: utils.StatusError{},
	}

	// it is an error to set this field in an http client request
	up.Request.RequestURI = ""

	return up
}

func (up *upstream) configureTLS(p *proxy) *upstream {
	// if previous error - continue
	if up.Error != (utils.StatusError{}) {
		return up
	}

	// get secret for this request - if any
	untypedSecret, err := spec.SecretStore.Get(p.Target.GetSSLCertificates(p.Source.Host).SecretName, p.Target.ObjectMeta.Namespace)
	if err != nil {
		up.Error = utils.StatusError{Code: http.StatusInternalServerError, Err: err}
		return up
	}

	tlsConfig := &tls.Config{}
	caCertPool := x509.NewCertPool()

	// ssl is not configured for this request
	if untypedSecret == nil {

		// if upstream option is being used, if the scheme
		// is https we need to add the root ca bundle
		if strings.Compare(up.Request.URL.Scheme, "https") != 0 {
			logrus.Debug("TLS not configured for this proxy")
			return up
		}

	} else {

		secret, ok := untypedSecret.(api.Secret)
		if !ok {
			up.Error = utils.StatusError{Code: http.StatusInternalServerError, Err: errors.New("the secret store is corrupted")}
			return up
		}

		// server side tls must be configured
		cert, err := spec.X509KeyPair(secret)
		if err != nil {
			up.Error = utils.StatusError{Code: http.StatusInternalServerError, Err: err}
			return up
		}
		tlsConfig.Certificates = []tls.Certificate{*cert}

		if secret.Data["tls.ca"] != nil {
			caCertPool.AppendCertsFromPEM(secret.Data["tls.ca"])
		}

		if viper.GetBool("disable-tls-cn-validation") {
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
	transport := &http.Transport{TLSClientConfig: tlsConfig}
	up.Client.Transport = transport

	return up

}

func (up *upstream) setUpstreamURL(p *proxy) *upstream {

	if up.Error != (utils.StatusError{}) {
		return up
	}

	u, err := p.setK8sDiscoveredURI()

	if err != nil {
		up.Error = *err
	} else {
		u.Path = utils.ComputeTargetPath(p.Target.Spec.Path, p.Target.Spec.Target, p.Source.URL.Path)
		u.RawPath = utils.ComputeTargetPath(p.Target.Spec.Path, p.Target.Spec.Target, p.Source.URL.EscapedPath())
		u.ForceQuery = p.Source.URL.ForceQuery
		u.RawQuery = p.Source.URL.RawQuery
		u.Fragment = p.Source.URL.Fragment

		up.Request.URL = u
	}

	return up

}

func (up *upstream) setUpstreamHeaders(p *proxy) *upstream {

	if up.Error != (utils.StatusError{}) {
		return up
	}

	// upstream request doesn't need the apikey
	// remove it.
	up.Request.Header.Del("apikey")

	up.Request.Header.Add("X-Forwarded-For", p.Source.RemoteAddr)

	return up

}

func (up *upstream) performProxy(trace opentracing.Span) *upstream {

	if up.Error != (utils.StatusError{}) {
		return up
	}

	logrus.Infof("upstream url: %s", up.Request.URL.String())

	err := trace.Tracer().Inject(trace.Context(),
		opentracing.TextMap,
		opentracing.HTTPHeadersCarrier(up.Request.Header))

	if err != nil {
		logrus.Error("could not inject headers")
	}

	resp, err := up.Client.Do(up.Request)
	if err != nil {
		up.Error = utils.StatusError{
			Code: http.StatusInternalServerError,
			Err:  err,
		}
	} else {
		up.Response = resp
	}

	return up

}

func (p *proxy) setK8sDiscoveredURI() (*url.URL, *utils.StatusError) {

	scheme := "http"

	if *p.Target.GetSSLCertificates(p.Source.Host) != (spec.SSL{}) {
		scheme = "https"
	}

	untypedSvc, err := spec.ServiceStore.Get(p.Target.Spec.Service, p.Source.Header)
	if err != nil || untypedSvc == nil {
		return nil, &utils.StatusError{Code: http.StatusNotFound, Err: errors.New("no matching services")}
	}

	svc, ok := untypedSvc.(spec.Service)
	if !ok {
		return nil, &utils.StatusError{Code: http.StatusNotFound, Err: errors.New("no matching services")}
	}

	uri := fmt.Sprintf("%s.%s.svc.cluster.local",
		svc.Name,
		p.Target.ObjectMeta.Namespace,
	)

	if viper.GetBool("enable-cluster-ip") {
		uri = svc.ClusterIP
	}

	return &url.URL{
		Scheme: scheme,
		Host: fmt.Sprintf("%s:%d",
			uri,
			p.Target.Spec.Service.Port,
		),
	}, nil

}
