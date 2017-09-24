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
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/northwesternmutual/kanali/config"
	"github.com/northwesternmutual/kanali/metrics"
	"github.com/northwesternmutual/kanali/spec"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/mocktracer"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"k8s.io/kubernetes/pkg/api"
)

type mockHTTPClient struct{}

func (cli *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	responseRecorder := &httptest.ResponseRecorder{}
	mockTracer := mocktracer.New()

	_, err := mockTracer.Extract(
		opentracing.TextMap,
		opentracing.HTTPHeadersCarrier(req.Header),
	)
	if err != nil {
		return nil, errors.New("error extracting headers")
	}

	if req.URL.Path == "/error" {
		return nil, errors.New("expected error")
	}

	responseRecorder.Code = http.StatusOK
	responseRecorder.Body = bytes.NewBuffer([]byte("response body"))
	return responseRecorder.Result(), nil
}

func TestProxyPassGetName(t *testing.T) {
	step := ProxyPassStep{}
	assert.Equal(t, step.GetName(), "Proxy Pass", "step name is incorrect")
}

func TestPreformTargetProxy(t *testing.T) {
	testMetrics := &metrics.Metrics{}
	mockTracer := mocktracer.New()
	testReqOne, _ := http.NewRequest("GET", "https://foo.bar.com/?foo=bar", bytes.NewReader([]byte("test data")))
	testReqTwo, _ := http.NewRequest("GET", "https://foo.bar.com/error", bytes.NewReader([]byte("test data")))

	testSpanOne := mockTracer.StartSpan("test span one")
	resp, err := preformTargetProxy(&mockHTTPClient{}, testReqOne, testMetrics, testSpanOne)
	testSpanOne.Finish()
	assert.Nil(t, err)
	assert.Equal(t, resp.StatusCode, 200)
	assert.Equal(t, (*testMetrics)[0].Name, "total_target_time")
	assert.False(t, (*testMetrics)[0].Index)

	testSpanTwo := mockTracer.StartSpan("test span two")
	_, err = preformTargetProxy(&mockHTTPClient{}, testReqTwo, testMetrics, testSpanTwo)
	testSpanTwo.Finish()
	assert.Equal(t, err.Error(), "expected error")
}

func TestGetTargetHost(t *testing.T) {
	spec.ServiceStore.Set(spec.Service{
		Name:      "bar",
		Namespace: "foo",
		ClusterIP: "1.2.3.4",
		Port:      8080,
	})
	req, _ := http.NewRequest("GET", "http://foo.bar.com", nil)

	proxyOne := &spec.APIProxy{
		ObjectMeta: api.ObjectMeta{
			Name:      "exampleAPIProxyOne",
			Namespace: "foo",
		},
		Spec: spec.APIProxySpec{
			Path: "/api/v1/accounts",
			Service: spec.Service{
				Name:      "bar",
				Namespace: "foo",
				Port:      8080,
			},
		},
	}

	proxyTwo := &spec.APIProxy{
		ObjectMeta: api.ObjectMeta{
			Name:      "exampleAPIProxyOne",
			Namespace: "foo",
		},
		Spec: spec.APIProxySpec{
			Path:   "/api/v1/accounts",
			Target: "/foo/bar",
			Service: spec.Service{
				Name:      "bar",
				Namespace: "foo",
				Port:      8080,
			},
			SSL: spec.SSL{
				SecretName: "mysecretname",
			},
		},
	}

	urlOne, _ := getTargetHost(proxyOne, req)
	assert.Equal(t, *urlOne, url.URL{
		Scheme: "http",
		Host:   "bar.foo.svc.cluster.local:8080",
	})

	viper.SetDefault(config.FlagProxyEnableClusterIP.GetLong(), true)

	urlOne, _ = getTargetHost(proxyOne, req)
	assert.Equal(t, *urlOne, url.URL{
		Scheme: "http",
		Host:   "1.2.3.4:8080",
	})

	viper.SetDefault(config.FlagProxyEnableClusterIP.GetLong(), false)

	urlTwo, _ := getTargetHost(proxyTwo, req)
	assert.Equal(t, *urlTwo, url.URL{
		Scheme: "https",
		Host:   "bar.foo.svc.cluster.local:8080",
	})

	viper.SetDefault(config.FlagProxyEnableClusterIP.GetLong(), true)

	urlTwo, _ = getTargetHost(proxyTwo, req)
	assert.Equal(t, *urlTwo, url.URL{
		Scheme: "https",
		Host:   "1.2.3.4:8080",
	})
}
