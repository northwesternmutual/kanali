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

package handlers

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/northwesternmutual/kanali/config"
	"github.com/northwesternmutual/kanali/logging"
	"github.com/northwesternmutual/kanali/metrics"
	"github.com/northwesternmutual/kanali/spec"
	"github.com/northwesternmutual/kanali/tracer"
	"github.com/northwesternmutual/kanali/utils"
	"github.com/opentracing/opentracing-go/mocktracer"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/pkg/api/v1"
)

func init() {
	logging.Init(nil)
}

func TestIncomingRequest(t *testing.T) {
	mock, _ := json.Marshal([]spec.Route{
		{
			Route:  "/foo",
			Code:   200,
			Method: "GET",
			Body: map[string]interface{}{
				"foo": "bar",
			},
		},
	})
	spec.MockResponseStore.Set(v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "mockConfigMap",
			Namespace: "foo",
		},
		Data: map[string]string{
			"response": string(mock),
		},
	})
	spec.ProxyStore.Set(spec.APIProxy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "exampleAPIProxyOne",
			Namespace: "foo",
		},
		Spec: spec.APIProxySpec{
			Path:   "/api/v1/accounts",
			Target: "/",
			Service: spec.Service{
				Name: "dummyService",
				Port: 8080,
			},
			Mock: &spec.Mock{
				ConfigMapName: "mockConfigMap",
			},
		},
	})
	metrics := &metrics.Metrics{}
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "http://foo.bar.com/api/v1/accounts/foo", nil)
	mockTracer := mocktracer.New()

	viper.SetDefault(config.FlagProxyEnableMockResponses.GetLong(), true)
	defer viper.Reset()
	span := mockTracer.StartSpan("test span")
	err := IncomingRequest(context.Background(), &spec.APIProxy{}, metrics, writer, request, span)
	span.Finish()
	assert.Nil(t, err)
	assert.Equal(t, mockTracer.FinishedSpans()[0].Tags()[tracer.KanaliProxyName], "exampleAPIProxyOne")
	assert.Equal(t, mockTracer.FinishedSpans()[0].Tags()[tracer.KanaliProxyNamespace], "foo")
	assert.Equal(t, (*metrics)[0].Name, "proxy_name")
	assert.Equal(t, (*metrics)[0].Value, "exampleAPIProxyOne")
	assert.Equal(t, (*metrics)[1].Name, "proxy_namespace")
	assert.Equal(t, (*metrics)[1].Value, "foo")
	assert.Equal(t, (*metrics)[2].Name, "http_response_code")
	assert.Equal(t, (*metrics)[2].Value, "200")

	response := writer.Result()
	assert.Equal(t, response.StatusCode, 200)
	assert.Equal(t, response.Header, http.Header{"Content-Type": []string{"application/json"}})
	body, _ := ioutil.ReadAll(response.Body)
	assert.Equal(t, string(body), `{"foo":"bar"}`)

	viper.SetDefault(config.FlagProxyEnableMockResponses.GetLong(), false)
	span = mockTracer.StartSpan("test span")
	defer span.Finish()
	err = IncomingRequest(context.Background(), &spec.APIProxy{}, metrics, writer, request, span)
	statusErr := err.(utils.Error)
	assert.Equal(t, statusErr.Status(), 404)
	assert.Equal(t, statusErr.Error(), "no matching services")
}

func TestMockIsDefined(t *testing.T) {
	spec.ProxyStore.Set(spec.APIProxy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "exampleAPIProxyOne",
			Namespace: "foo",
		},
		Spec: spec.APIProxySpec{
			Path:   "/api/v1/accounts",
			Target: "/",
			Service: spec.Service{
				Name: "dummyService",
				Port: 8080,
			},
			Mock: &spec.Mock{
				ConfigMapName: "mockConfigMap",
			},
		},
	})
	spec.ProxyStore.Set(spec.APIProxy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "exampleAPIProxyTwo",
			Namespace: "foo",
		},
		Spec: spec.APIProxySpec{
			Path:   "/api/v1/properties",
			Target: "/",
			Service: spec.Service{
				Name: "dummyService",
				Port: 8080,
			},
		},
	})

	result := mockIsDefined("/api/v1/accounts/foo")
	assert.True(t, result)

	result = mockIsDefined("/api/v1/clients/foo")
	assert.False(t, result)

	result = mockIsDefined("/api/v1/properties/foo")
	assert.False(t, result)

}
