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
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/northwesternmutual/kanali/metrics"
	"github.com/northwesternmutual/kanali/spec"
	"github.com/northwesternmutual/kanali/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/stretchr/testify/assert"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/unversioned"
)

func TestValidateProxyGetName(t *testing.T) {
	assert := assert.New(t)
	step := ValidateProxyStep{}
	assert.Equal(step.GetName(), "Validate Proxy", "step name is incorrect")
}

func TestValidateProxy(t *testing.T) {
	assert := assert.New(t)
	step := ValidateProxyStep{}

	proxyStore := spec.ProxyStore
	proxyList := getTestAPIProxyListForValidateProxy()
	proxyStore.Clear()
	proxyStore.Set(proxyList.Proxies[0])
	proxyStore.Set(proxyList.Proxies[1])

	urlOne, _ := url.Parse("https://www.foo.bar.com/api/v1/accounts/one/two")
	urlTwo, _ := url.Parse("https://www.foo.bar.com/api/v1/field/one/two")
	urlThree, _ := url.Parse("https://www.foo.bar.com/")
	urlFour, _ := url.Parse("https://www.foo.bar.com/foo/bar")

	assert.Nil(step.Do(context.Background(), nil, &metrics.Metrics{}, nil, &http.Request{URL: urlOne}, nil, opentracing.StartSpan("test span")), "expected proxy to be found")
	assert.Nil(step.Do(context.Background(), nil, &metrics.Metrics{}, nil, &http.Request{URL: urlTwo}, nil, opentracing.StartSpan("test span")), "expected proxy to be found")
	assert.Equal(utils.StatusError{Code: http.StatusNotFound, Err: errors.New("proxy not found")}, step.Do(context.Background(), nil, &metrics.Metrics{}, nil, &http.Request{URL: urlThree}, nil, opentracing.StartSpan("test span")), "expected proxy to not exist")
	assert.Equal(utils.StatusError{Code: http.StatusNotFound, Err: errors.New("proxy not found")}, step.Do(context.Background(), nil, &metrics.Metrics{}, nil, &http.Request{URL: urlFour}, nil, opentracing.StartSpan("test span")), "expected proxy to not exist")
}

func getTestAPIProxyListForValidateProxy() *spec.APIProxyList {

	return &spec.APIProxyList{
		TypeMeta: unversioned.TypeMeta{},
		ListMeta: unversioned.ListMeta{},
		Proxies: []spec.APIProxy{
			{
				TypeMeta: unversioned.TypeMeta{},
				ObjectMeta: api.ObjectMeta{
					Name:      "exampleAPIProxyOne",
					Namespace: "foo",
				},
				Spec: spec.APIProxySpec{
					Path: "api/v1/accounts",
				},
			},
			{
				TypeMeta: unversioned.TypeMeta{},
				ObjectMeta: api.ObjectMeta{
					Name:      "exampleAPIProxyTwo",
					Namespace: "foo",
				},
				Spec: spec.APIProxySpec{
					Path: "/api/v1/field",
				},
			},
		},
	}

}
