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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes/fake"

	"github.com/northwesternmutual/kanali/pkg/apis/kanali.io/v2"
	"github.com/northwesternmutual/kanali/pkg/log"
	"github.com/northwesternmutual/kanali/test/builder"
)

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
				originalReq: builder.NewHTTPRequest().NewOrDie(),
				upstreamReq: builder.NewHTTPRequest().WithHost("http://foo.bar.com").NewOrDie(),
			},
		},
		{
			config: builder.NewTLSConfigBuilder().WithSystemRoots().WithInsecure().WithVerify().NewOrDie(),
			err:    false,
			step: proxyPassStep{
				originalReq: builder.NewHTTPRequest().NewOrDie(),
				upstreamReq: builder.NewHTTPRequest().WithHost("https://foo.bar.com").NewOrDie(),
				proxy:       builder.NewApiProxy("foo", "bar").NewOrDie(),
			},
		},
		{
			config: builder.NewTLSConfigBuilder().WithSystemRoots().WithInsecure().WithVerify().NewOrDie(),
			err:    false,
			step: proxyPassStep{
				originalReq: builder.NewHTTPRequest().NewOrDie(),
				upstreamReq: builder.NewHTTPRequest().WithHost("https://foo.bar.com").NewOrDie(),
				proxy:       builder.NewApiProxy("foo", "bar").WithSecret("").NewOrDie(),
			},
		},
		{
			config: nil,
			err:    true,
			step: proxyPassStep{
				v1Interface: informers.NewSharedInformerFactory(fake.NewSimpleClientset(), 500*time.Millisecond).Core().V1(),
				originalReq: builder.NewHTTPRequest().NewOrDie(),
				upstreamReq: builder.NewHTTPRequest().WithHost("https://foo.bar.com").NewOrDie(),
				proxy:       builder.NewApiProxy("foo", "bar").WithSecret("foo").NewOrDie(),
			},
		},
		{
			config: nil,
			err:    true,
			step: proxyPassStep{
				v1Interface: informers.NewSharedInformerFactory(fake.NewSimpleClientset(), 500*time.Millisecond).Core().V1(),
				originalReq: builder.NewHTTPRequest().NewOrDie(),
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
				v1Interface: informers.NewSharedInformerFactory(fake.NewSimpleClientset(), 500*time.Millisecond).Core().V1(),
				originalReq: builder.NewHTTPRequest().NewOrDie(),
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
				v1Interface: informers.NewSharedInformerFactory(fake.NewSimpleClientset(), 500*time.Millisecond).Core().V1(),
				originalReq: builder.NewHTTPRequest().NewOrDie(),
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
				v1Interface: informers.NewSharedInformerFactory(fake.NewSimpleClientset(), 500*time.Millisecond).Core().V1(),
				originalReq: builder.NewHTTPRequest().NewOrDie(),
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
			config: builder.NewTLSConfigBuilder().WithSystemRoots().WithKeyPair(tlsAssets.ServerCert, tlsAssets.ServerKey).WithInsecure().WithVerify().NewOrDie(),
			err:    false,
			step: proxyPassStep{
				v1Interface: informers.NewSharedInformerFactory(fake.NewSimpleClientset(), 500*time.Millisecond).Core().V1(),
				originalReq: builder.NewHTTPRequest().NewOrDie(),
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
			config: builder.NewTLSConfigBuilder().WithSystemRoots().WithKeyPair(tlsAssets.ServerCert, tlsAssets.ServerKey).WithInsecure().WithVerify().NewOrDie(),
			err:    false,
			step: proxyPassStep{
				v1Interface: informers.NewSharedInformerFactory(fake.NewSimpleClientset(), 500*time.Millisecond).Core().V1(),
				originalReq: builder.NewHTTPRequest().NewOrDie(),
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
			config: builder.NewTLSConfigBuilder().WithSystemRoots().WithCustomCA(tlsAssets.CACert).WithInsecure().WithVerify().NewOrDie(),
			err:    false,
			step: proxyPassStep{
				v1Interface: informers.NewSharedInformerFactory(fake.NewSimpleClientset(), 500*time.Millisecond).Core().V1(),
				originalReq: builder.NewHTTPRequest().NewOrDie(),
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
			config: builder.NewTLSConfigBuilder().WithSystemRoots().WithCustomCA(tlsAssets.CACert).WithInsecure().WithVerify().NewOrDie(),
			err:    false,
			step: proxyPassStep{
				v1Interface: informers.NewSharedInformerFactory(fake.NewSimpleClientset(), 500*time.Millisecond).Core().V1(),
				originalReq: builder.NewHTTPRequest().NewOrDie(),
				upstreamReq: builder.NewHTTPRequest().WithHost("https://foo.bar.com").NewOrDie(),
				proxy:       builder.NewApiProxy("foo", "bar").WithSecret("foo").NewOrDie(),
			},
			prep: func(step proxyPassStep) {
				step.v1Interface.Secrets().Informer().GetStore().Add(
					builder.NewSecretBuilder("foo", "bar").WithAnnotation("kanali.io/enabled", "true").WithAnnotation("kanali.io/ca", "foo.bar").WithData("foo.bar", tlsAssets.CACert).NewOrDie(),
				)
			},
		},
	}

	for _, test := range tests {
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
