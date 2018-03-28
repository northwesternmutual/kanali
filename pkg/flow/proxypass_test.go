// Copyright (c) 2018 Northwestern Mutual.
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
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes/fake"

	"github.com/northwesternmutual/kanali/cmd/kanali/app/options"
	"github.com/northwesternmutual/kanali/pkg/apis/kanali.io/v2"
	"github.com/northwesternmutual/kanali/pkg/log"
	"github.com/northwesternmutual/kanali/test/builder"
)

func TestProxyPassName(t *testing.T) {
	assert.Equal(t, "Proxy Pass", ProxyPassStep(nil).Name())
}

func TestGetServiceLabelSet(t *testing.T) {
	tests := []struct {
		proxy    *v2.ApiProxy
		headers  http.Header
		labels   labels.Set
		defaults map[string]string
	}{
		{
			proxy: builder.NewApiProxy("foo", "bar").WithTargetBackendDynamicService(8080,
				v2.Label{Name: "foo", Value: "bar"},
			).NewOrDie(),
			headers: nil,
			labels: map[string]string{
				"foo": "bar",
			},
			defaults: make(map[string]string),
		},
		{
			proxy:    builder.NewApiProxy("foo", "bar").WithTargetBackendDynamicService(8080).NewOrDie(),
			headers:  nil,
			labels:   map[string]string{},
			defaults: make(map[string]string),
		},
		{
			proxy: builder.NewApiProxy("foo", "bar").WithTargetBackendDynamicService(8080,
				v2.Label{Name: "foo", Value: "bar"},
				v2.Label{Name: "bar", Header: "foo"},
			).NewOrDie(),
			headers: make(http.Header),
			labels: map[string]string{
				"foo": "bar",
				"bar": "",
			},
			defaults: make(map[string]string),
		},
		{
			proxy: builder.NewApiProxy("foo", "bar").WithTargetBackendDynamicService(8080,
				v2.Label{Name: "foo", Value: "bar"},
				v2.Label{Name: "bar", Header: "foo"},
			).NewOrDie(),
			headers: http.Header(map[string][]string{
				"Foo": {"car"},
			}),
			labels: map[string]string{
				"foo": "bar",
				"bar": "car",
			},
			defaults: make(map[string]string),
		},
		{
			proxy: builder.NewApiProxy("foo", "bar").WithTargetBackendDynamicService(8080,
				v2.Label{Name: "bar", Header: "foo"},
			).NewOrDie(),
			headers: http.Header(map[string][]string{
				"Foo": {"car"},
			}),
			labels: map[string]string{
				"bar": "car",
			},
			defaults: make(map[string]string),
		},
		{
			proxy: builder.NewApiProxy("foo", "bar").WithTargetBackendDynamicService(8080,
				v2.Label{Name: "foo", Value: "bar"},
				v2.Label{Name: "bar", Header: "foo"},
			).NewOrDie(),
			headers: make(http.Header),
			labels: map[string]string{
				"foo": "bar",
				"bar": "default",
			},
			defaults: map[string]string{
				"foo": "default",
			},
		},
		{
			proxy: builder.NewApiProxy("foo", "bar").WithTargetBackendDynamicService(8080,
				v2.Label{Name: "foo", Value: "bar"},
				v2.Label{Name: "bar", Header: "foo"},
			).NewOrDie(),
			headers: http.Header(map[string][]string{
				"Foo": {"car"},
			}),
			labels: map[string]string{
				"foo": "bar",
				"bar": "car",
			},
			defaults: map[string]string{
				"foo": "default",
			},
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.labels, getServiceLabelSet(test.proxy, test.headers, test.defaults))
	}
}

func BenchmarkGetServiceLabelSet(b *testing.B) {
	proxy := builder.NewApiProxy("foo", "bar").WithTargetBackendDynamicService(8080,
		v2.Label{Name: "foo", Value: "bar"},
		v2.Label{Name: "bar", Header: "foo"},
	).NewOrDie()
	headers := make(http.Header)
	defaults := map[string]string{
		"foo": "default",
	}

	for n := 0; n < b.N; n++ {
		getServiceLabelSet(proxy, headers, defaults)
	}
}

func TestCopyHeader(t *testing.T) {
	original := http.Header(map[string][]string{
		"Foo": {"bar"},
	})
	copy := make(http.Header)
	assert.Equal(t, 0, len(copy))
	copyHeader(copy, original)
	assert.Equal(t, 1, len(copy))
}

func TestConfigureTLS(t *testing.T) {
	core, _ := observer.New(zap.NewAtomicLevelAt(zapcore.DebugLevel))
	defer log.SetLogger(zap.New(core)).Restore()

	tlsAssets := builder.NewTLSBuilder(nil, nil).NewOrDie()

	i := informers.NewSharedInformerFactory(fake.NewSimpleClientset(), 500*time.Millisecond).Core().V1()

	tests := []struct {
		config *tls.Config
		err    bool
		step   proxyPassStep
		prep   func(proxyPassStep)
	}{
		{
			config: nil,
			err:    false,
			step: proxyPassStep{
				upstreamReq: builder.NewHTTPRequest().WithHost("http://foo.bar.com").NewOrDie(),
			},
		},
		{
			config: builder.NewTLSConfigBuilder().WithSystemRoots().WithInsecure().WithVerify().NewOrDie(),
			err:    false,
			step: proxyPassStep{
				upstreamReq: builder.NewHTTPRequest().WithHost("https://foo.bar.com").NewOrDie(),
				proxy:       builder.NewApiProxy("foo", "bar").NewOrDie(),
			},
		},
		{
			config: builder.NewTLSConfigBuilder().WithSystemRoots().WithInsecure().WithVerify().NewOrDie(),
			err:    false,
			step: proxyPassStep{
				upstreamReq: builder.NewHTTPRequest().WithHost("https://foo.bar.com").NewOrDie(),
				proxy:       builder.NewApiProxy("foo", "bar").WithSecret("").NewOrDie(),
			},
		},
		{
			config: nil,
			err:    true,
			step: proxyPassStep{
				upstreamReq: builder.NewHTTPRequest().WithHost("https://foo.bar.com").NewOrDie(),
				proxy:       builder.NewApiProxy("foo", "bar").WithSecret("foo").NewOrDie(),
			},
		},
		{
			config: nil,
			err:    true,
			step: proxyPassStep{
				upstreamReq: builder.NewHTTPRequest().WithHost("https://foo.bar.com").NewOrDie(),
				proxy:       builder.NewApiProxy("foo", "bar").WithSecret("foo").NewOrDie(),
			},
			prep: func(step proxyPassStep) {
				step.v1Interface.Secrets().Informer().GetStore().Add(
					builder.NewSecretBuilder("foo", "bar").NewOrDie(),
				)
			},
		},
		{
			config: nil,
			err:    true,
			step: proxyPassStep{
				upstreamReq: builder.NewHTTPRequest().WithHost("https://foo.bar.com").NewOrDie(),
				proxy:       builder.NewApiProxy("foo", "bar").WithSecret("foo").NewOrDie(),
			},
			prep: func(step proxyPassStep) {
				step.v1Interface.Secrets().Informer().GetStore().Add(
					builder.NewSecretBuilder("foo", "bar").WithAnnotation("kanali.io/enabled", "true").NewOrDie(),
				)
			},
		},
		{
			config: nil,
			err:    true,
			step: proxyPassStep{
				upstreamReq: builder.NewHTTPRequest().WithHost("https://foo.bar.com").NewOrDie(),
				proxy:       builder.NewApiProxy("foo", "bar").WithSecret("foo").NewOrDie(),
			},
			prep: func(step proxyPassStep) {
				step.v1Interface.Secrets().Informer().GetStore().Add(
					builder.NewSecretBuilder("foo", "bar").WithAnnotation("kanali.io/enabled", "true").WithData("tls.key", []byte("foo")).WithData("tls.crt", []byte("bar")).NewOrDie(),
				)
			},
		},
		{
			config: nil,
			err:    true,
			step: proxyPassStep{
				upstreamReq: builder.NewHTTPRequest().WithHost("https://foo.bar.com").NewOrDie(),
				proxy:       builder.NewApiProxy("foo", "bar").WithSecret("foo").NewOrDie(),
			},
			prep: func(step proxyPassStep) {
				step.v1Interface.Secrets().Informer().GetStore().Add(
					builder.NewSecretBuilder("foo", "bar").WithAnnotation("kanali.io/enabled", "true").WithData("tls.key", []byte("foo")).WithData("tls.ca", []byte("bar")).NewOrDie(),
				)
			},
		},
		{
			config: builder.NewTLSConfigBuilder().WithKeyPair(tlsAssets.ServerCert, tlsAssets.ServerKey).WithInsecure().WithVerify().NewOrDie(),
			err:    false,
			step: proxyPassStep{
				upstreamReq: builder.NewHTTPRequest().WithHost("https://foo.bar.com").NewOrDie(),
				proxy:       builder.NewApiProxy("foo", "bar").WithSecret("foo").NewOrDie(),
			},
			prep: func(step proxyPassStep) {
				step.v1Interface.Secrets().Informer().GetStore().Add(
					builder.NewSecretBuilder("foo", "bar").WithAnnotation("kanali.io/enabled", "true").WithData("tls.key", tlsAssets.ServerKey).WithData("tls.crt", tlsAssets.ServerCert).NewOrDie(),
				)
			},
		},
		{
			config: builder.NewTLSConfigBuilder().WithKeyPair(tlsAssets.ServerCert, tlsAssets.ServerKey).WithInsecure().WithVerify().NewOrDie(),
			err:    false,
			step: proxyPassStep{
				upstreamReq: builder.NewHTTPRequest().WithHost("https://foo.bar.com").NewOrDie(),
				proxy:       builder.NewApiProxy("foo", "bar").WithSecret("foo").NewOrDie(),
			},
			prep: func(step proxyPassStep) {
				step.v1Interface.Secrets().Informer().GetStore().Add(
					builder.NewSecretBuilder("foo", "bar").WithAnnotation("kanali.io/enabled", "true").WithAnnotation("kanali.io/key", "foo.bar").WithAnnotation("kanali.io/cert", "bar.foo").WithData("foo.bar", tlsAssets.ServerKey).WithData("bar.foo", tlsAssets.ServerCert).NewOrDie(),
				)
			},
		},
		{
			config: builder.NewTLSConfigBuilder().WithCustomCA(tlsAssets.CACert).WithInsecure().WithVerify().NewOrDie(),
			err:    false,
			step: proxyPassStep{
				upstreamReq: builder.NewHTTPRequest().WithHost("https://foo.bar.com").NewOrDie(),
				proxy:       builder.NewApiProxy("foo", "bar").WithSecret("foo").NewOrDie(),
			},
			prep: func(step proxyPassStep) {
				step.v1Interface.Secrets().Informer().GetStore().Add(
					builder.NewSecretBuilder("foo", "bar").WithAnnotation("kanali.io/enabled", "true").WithData("tls.ca", tlsAssets.CACert).NewOrDie(),
				)
			},
		},
		{
			config: builder.NewTLSConfigBuilder().WithCustomCA(tlsAssets.CACert).WithInsecure().WithVerify().NewOrDie(),
			err:    false,
			step: proxyPassStep{
				upstreamReq: builder.NewHTTPRequest().WithHost("https://foo.bar.com").NewOrDie(),
				proxy:       builder.NewApiProxy("foo", "bar").WithSecret("foo").NewOrDie(),
			},
			prep: func(step proxyPassStep) {
				step.v1Interface.Secrets().Informer().GetStore().Add(
					builder.NewSecretBuilder("foo", "bar").WithAnnotation("kanali.io/enabled", "true").WithAnnotation("kanali.io/ca", "foo.bar").WithData("foo.bar", tlsAssets.CACert).NewOrDie(),
				)
			},
		},
		{
			config: builder.NewTLSConfigBuilder().WithSystemRoots().WithInsecure().WithVerify().NewOrDie(),
			err:    false,
			step: proxyPassStep{
				upstreamReq: builder.NewHTTPRequest().WithHost("https://foo.bar.com").NewOrDie(),
				proxy:       builder.NewApiProxy("foo", "bar").NewOrDie(),
			},
		},
	}

	for _, test := range tests {
		test.step.v1Interface = i
		test.step.originalReq = builder.NewHTTPRequest().NewOrDie()

		if test.prep != nil {
			test.prep(test.step)
		}
		cfg, err := test.step.configureTLS()
		if !test.err {
			assert.Nil(t, err)
			if cfg != nil {
				assert.Equal(t, test.config.RootCAs, cfg.RootCAs)
				assert.Equal(t, test.config.InsecureSkipVerify, cfg.InsecureSkipVerify)
				assert.Equal(t, test.config.Certificates, cfg.Certificates)
				if test.config.VerifyPeerCertificate != nil {
					assert.NotNil(t, cfg.VerifyPeerCertificate)
				} else {
					assert.Nil(t, cfg.VerifyPeerCertificate)
				}
			} else {
				assert.Nil(t, test.config)
			}
		} else {
			assert.NotNil(t, err)
		}
	}
}

func BenchmarkConfigureTLS(b *testing.B) {
	tlsAssets := builder.NewTLSBuilder(nil, nil).NewOrDie()
	step := proxyPassStep{
		v1Interface: informers.NewSharedInformerFactory(fake.NewSimpleClientset(), 500*time.Millisecond).Core().V1(),
		originalReq: builder.NewHTTPRequest().NewOrDie(),
		upstreamReq: builder.NewHTTPRequest().WithHost("https://foo.bar.com").NewOrDie(),
		proxy:       builder.NewApiProxy("foo", "bar").WithSecret("foo").NewOrDie(),
	}
	step.v1Interface.Secrets().Informer().GetStore().Add(
		builder.NewSecretBuilder("foo", "bar").WithAnnotation("kanali.io/enabled", "true").WithAnnotation("kanali.io/ca", "foo.bar").WithData("foo.bar", tlsAssets.CACert).NewOrDie(),
	)
	for i := 0; i < b.N; i++ {
		step.configureTLS()
	}
}

func TestUserDefinedSSL(t *testing.T) {
	assert.True(t, (proxyPassStep{
		proxy: builder.NewApiProxy("foo", "bar").WithSecret("foo").NewOrDie(),
	}).userDefinedSSL())
	assert.False(t, (proxyPassStep{
		proxy: builder.NewApiProxy("foo", "bar").NewOrDie(),
	}).userDefinedSSL())
}

func TestServiceDetails(t *testing.T) {
	core, _ := observer.New(zap.NewAtomicLevelAt(zapcore.DebugLevel))
	defer log.SetLogger(zap.New(core)).Restore()

	i := informers.NewSharedInformerFactory(fake.NewSimpleClientset(), 500*time.Millisecond).Core().V1()

	tests := []struct {
		expectedScheme, expectedHost string
		step                         proxyPassStep
		pre                          func(proxyPassStep)
		expectedErr                  bool
	}{
		{
			expectedErr: true,
			step: proxyPassStep{
				proxy: builder.NewApiProxy("foo", "bar").WithTargetBackendDynamicService(8080, v2.Label{
					Name:  "foo",
					Value: "bar",
				}).NewOrDie(),
			},
		},
		{
			expectedErr: true,
			step: proxyPassStep{
				proxy: builder.NewApiProxy("foo", "bar").NewOrDie(),
			},
		},
		{
			expectedErr: true,
			step: proxyPassStep{
				proxy: builder.NewApiProxy("foo", "bar").WithSecret("foo").WithTargetBackendStaticService("foo", 8080).NewOrDie(),
			},
		},
		{
			expectedErr:    false,
			expectedScheme: "https",
			expectedHost:   "1.2.3.4:8080",
			step: proxyPassStep{
				proxy: builder.NewApiProxy("foo", "bar").WithSecret("foo").WithTargetBackendStaticService("foo", 8080).NewOrDie(),
			},
			pre: func(step proxyPassStep) {
				viper.SetDefault(options.FlagProxyEnableClusterIP.GetLong(), true)
				step.v1Interface.Services().Informer().GetStore().Add(
					builder.NewServiceBuilder("foo", "bar").WithClusterIP("1.2.3.4").NewOrDie(),
				)
			},
		},
		{
			expectedErr:    false,
			expectedScheme: "http",
			expectedHost:   "foo.bar.svc.cluster.local:8080",
			step: proxyPassStep{
				proxy: builder.NewApiProxy("foo", "bar").WithTargetBackendStaticService("foo", 8080).NewOrDie(),
			},
			pre: func(step proxyPassStep) {
				step.v1Interface.Services().Informer().GetStore().Add(
					builder.NewServiceBuilder("foo", "bar").NewOrDie(),
				)
			},
		},
		{
			expectedErr:    false,
			expectedScheme: "http",
			expectedHost:   "1.2.3.4:8080",
			step: proxyPassStep{
				proxy: builder.NewApiProxy("foo", "bar").WithTargetBackendDynamicService(8080, v2.Label{
					Name:  "foo",
					Value: "bar",
				}).NewOrDie(),
			},
			pre: func(step proxyPassStep) {
				viper.SetDefault(options.FlagProxyEnableClusterIP.GetLong(), true)
				step.v1Interface.Services().Informer().GetStore().Add(
					builder.NewServiceBuilder("foo", "bar").WithClusterIP("1.2.3.4").WithLabel("foo", "bar").NewOrDie(),
				)
				step.v1Interface.Services().Informer().GetStore().Add(
					builder.NewServiceBuilder("car", "bar").WithClusterIP("1.2.3.4").WithLabel("foo", "bar").NewOrDie(),
				)
			},
		},
	}

	for _, test := range tests {
		test.step.v1Interface = i
		test.step.originalRespWriter = httptest.NewRecorder()
		test.step.originalReq = builder.NewHTTPRequest().NewOrDie()

		if test.pre != nil {
			test.pre(test.step)
		}
		scheme, host, err := test.step.serviceDetails()
		if test.expectedErr {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
			assert.Equal(t, test.expectedScheme, scheme)
			assert.Equal(t, test.expectedHost, host)
		}
		viper.Reset()
	}
}

func TestSetUpstreamURL(t *testing.T) {
	core, _ := observer.New(zap.NewAtomicLevelAt(zapcore.DebugLevel))
	defer log.SetLogger(zap.New(core)).Restore()

	i := informers.NewSharedInformerFactory(fake.NewSimpleClientset(), 500*time.Millisecond).Core().V1()

	tests := []struct {
		expectedURL *url.URL
		step        proxyPassStep
		pre         func(proxyPassStep)
		expectedErr bool
	}{
		{
			expectedErr: true,
			step: proxyPassStep{
				originalReq: builder.NewHTTPRequest().NewOrDie(),
				proxy:       builder.NewApiProxy("foo", "bar").WithTargetBackendStaticService("foo", 8080).NewOrDie(),
			},
		},
		{
			expectedErr: false,
			expectedURL: &url.URL{
				Scheme: "http",
				Host:   "foo.bar.svc.cluster.local:8080",
			},
			step: proxyPassStep{
				originalReq: builder.NewHTTPRequest().NewOrDie(),
				proxy:       builder.NewApiProxy("foo", "bar").WithTargetBackendStaticService("foo", 8080).NewOrDie(),
			},
			pre: func(step proxyPassStep) {
				step.v1Interface.Services().Informer().GetStore().Add(
					builder.NewServiceBuilder("foo", "bar").NewOrDie(),
				)
			},
		},
		{
			expectedErr: false,
			expectedURL: &url.URL{
				Scheme:  "http",
				Host:    "foo.bar.svc.cluster.local:8080",
				Path:    "/Go/",
				RawPath: "/%47%6f%2f",
			},
			step: proxyPassStep{
				originalReq: builder.NewHTTPRequest().WithPath("/%47%6f%2f").NewOrDie(),
				proxy:       builder.NewApiProxy("foo", "bar").WithTargetBackendStaticService("foo", 8080).NewOrDie(),
			},
			pre: func(step proxyPassStep) {
				step.v1Interface.Services().Informer().GetStore().Add(
					builder.NewServiceBuilder("foo", "bar").NewOrDie(),
				)
			},
		},
		{
			expectedErr: true,
			step: proxyPassStep{
				originalReq: builder.NewHTTPRequest().NewOrDie(),
				proxy:       builder.NewApiProxy("foo", "bar").WithTargetBackendEndpoint("foo://car.baz").NewOrDie(),
			},
		},
		{
			expectedErr: false,
			expectedURL: &url.URL{
				Scheme:  "http",
				Host:    "car.baz",
				Path:    "/Go/",
				RawPath: "/%47%6f%2f",
			},
			step: proxyPassStep{
				originalReq: builder.NewHTTPRequest().WithPath("/%47%6f%2f").NewOrDie(),
				proxy:       builder.NewApiProxy("foo", "bar").WithTargetBackendEndpoint("http://car.baz").NewOrDie(),
			},
		},
	}

	for _, test := range tests {
		test.step.v1Interface = i
		test.step.upstreamReq = &http.Request{
			URL: &url.URL{},
		}

		if test.pre != nil {
			test.pre(test.step)
		}
		err := test.step.setUpstreamURL()
		if test.expectedErr {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
			assert.Equal(t, test.expectedURL, test.step.upstreamReq.URL)
		}
		viper.Reset()
	}
}
