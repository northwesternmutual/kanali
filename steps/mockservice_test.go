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
	"testing"
  "net/http"
  "io/ioutil"
  "encoding/json"

  "k8s.io/client-go/pkg/api/v1"
	"github.com/stretchr/testify/assert"
  "github.com/northwesternmutual/kanali/spec"
  metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
  "github.com/northwesternmutual/kanali/metrics"
  opentracing "github.com/opentracing/opentracing-go"
)

func TestMockServiceGetName(t *testing.T) {
	step := MockServiceStep{}
	assert.Equal(t, step.GetName(), "Mock Service", "step name is incorrect")
}

func TestMockServiceDo(t *testing.T) {
  cms := getTestConfigMaps()
  spec.MockResponseStore.Clear()
  spec.MockResponseStore.Set(cms[0])
  balanceProxy := &spec.APIProxy{
    ObjectMeta: metav1.ObjectMeta{
      Name:      "proxy-one",
      Namespace: "foo",
    },
    Spec: spec.APIProxySpec{
      Path: "api/v1/balance",
      Target: "/car",
      Mock: &spec.Mock{
        ConfigMapName: "cm-one",
      },
    },
  }
  addressProxy := &spec.APIProxy{
    ObjectMeta: metav1.ObjectMeta{
      Name:      "proxy-one",
      Namespace: "foo",
    },
    Spec: spec.APIProxySpec{
      Path: "api/v1/address",
      Target: "/car",
      Mock: &spec.Mock{
        ConfigMapName: "cm-two",
      },
    },
  }
  accountsProxy := &spec.APIProxy{
    ObjectMeta: metav1.ObjectMeta{
      Name:      "proxy-one",
      Namespace: "foo",
    },
    Spec: spec.APIProxySpec{
      Path: "api/v1/accounts",
      Target: "/foo",
      Mock: &spec.Mock{
        ConfigMapName: "cm-one",
      },
    },
  }
  step := MockServiceStep{}

  m := &metrics.Metrics{}
  req, _ := http.NewRequest("GET", "http://foo.bar.com/api/v1/accounts", nil)
  res := &http.Response{}
  span := opentracing.StartSpan("test span")
  err := step.Do(context.Background(), accountsProxy, m, nil, req, res, span)
  assert.Nil(t, err)
  assert.Equal(t, m.Get("http_response_code").Value, "200")
  assert.Equal(t, res.Header.Get("Content-Type"), "application/json")
  body, _ := ioutil.ReadAll(res.Body)
  assert.Equal(t, string(body), `{"foo":"bar"}`)
  req, _ = http.NewRequest("GET", "http://foo.bar.com/api/v1/accounts/bar", nil)
  assert.Nil(t, step.Do(context.Background(), accountsProxy, m, nil, req, res, span))
  req, _ = http.NewRequest("GET", "http://foo.bar.com/api/v1/balance", nil)
  assert.Equal(t, step.Do(context.Background(), balanceProxy, m, nil, req, res, span).Error(), "no mock response found")
  req, _ = http.NewRequest("GET", "http://foo.bar.com/api/v1/address", nil)
  assert.Equal(t, step.Do(context.Background(), addressProxy, m, nil, req, res, span).Error(), "no mock response found")
}

func getTestConfigMaps() []v1.ConfigMap {

	mockOne, _ := json.Marshal([]spec.Route{
		{
			Route:  "/foo",
			Code:   200,
			Method: "GET",
			Body:   map[string]interface{}{
        "foo": "bar",
      },
		},
	})

	return []v1.ConfigMap{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "cm-one",
				Namespace: "foo",
			},
			Data: map[string]string{
				"response": string(mockOne),
			},
		},
	}

}
