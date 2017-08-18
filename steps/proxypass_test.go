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
	"net/http"
	"net/url"
	"testing"

	"github.com/northwesternmutual/kanali/spec"
	"github.com/northwesternmutual/kanali/utils"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/unversioned"
)

func TestProxyPassGetName(t *testing.T) {
	assert := assert.New(t)
	step := ProxyPassStep{}
	assert.Equal(step.GetName(), "Proxy Pass", "step name is incorrect")
}

func TestCreate(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(create(getTestInternalProxies()[0]), &upstream{
		Request: getTestInternalProxies()[0].Source,
		Client:  &http.Client{},
		Error:   utils.StatusError{},
	}, "object not what expected")

}

func TestSetUpstreamURL(t *testing.T) {
	spec.ServiceStore.Set(spec.Service{
		Name:      "my-service",
		Namespace: "foo",
		ClusterIP: "1.2.3.4",
		Port:      8080,
	})

	p := getTestInternalProxies()[3]
	result := create(p).setUpstreamURL(p)
	assert.Equal(t, result.Request.URL.Path, "/foo/bar/https://foo.bar.com")
	assert.Equal(t, result.Request.URL.EscapedPath(), "/foo/bar/https%3A%2F%2Ffoo.bar.com")
}

func TestSetK8sDiscoveredURI(t *testing.T) {
	assert := assert.New(t)

	spec.ServiceStore.Set(spec.Service{
		Name:      "bar",
		Namespace: "foo",
		ClusterIP: "1.2.3.4",
		Port:      8080,
	})

	proxies := getTestInternalProxies()

	urlOne, _ := proxies[0].setK8sDiscoveredURI()
	assert.Equal(*urlOne, url.URL{
		Scheme: "http",
		Host:   "bar.foo.svc.cluster.local:8080",
	})

	viper.SetDefault("enable-cluster-ip", true)

	urlOne, _ = proxies[0].setK8sDiscoveredURI()
	assert.Equal(*urlOne, url.URL{
		Scheme: "http",
		Host:   "1.2.3.4:8080",
	})

	viper.SetDefault("enable-cluster-ip", false)

	urlTwo, _ := proxies[1].setK8sDiscoveredURI()
	assert.Equal(*urlTwo, url.URL{
		Scheme: "https",
		Host:   "bar.foo.svc.cluster.local:8080",
	})

	viper.SetDefault("enable-cluster-ip", true)

	urlTwo, _ = proxies[1].setK8sDiscoveredURI()
	assert.Equal(*urlTwo, url.URL{
		Scheme: "https",
		Host:   "1.2.3.4:8080",
	})
}

func getTestInternalProxies() []*proxy {

	one, _ := url.Parse("http://www.foo.bar.com/api/v1/accounts")
	two, _ := url.Parse("http://www.foo.bar.com/api/v1/accounts/one/two")
	three, _ := url.Parse("http://www.foo.bar.com/api/v1/accounts/https%3A%2F%2Ffoo.bar.com")

	return []*proxy{
		{
			Source: &http.Request{
				URL: one,
			},
			Target: spec.APIProxy{
				TypeMeta: unversioned.TypeMeta{},
				ObjectMeta: api.ObjectMeta{
					Name:      "exampleAPIProxyOne",
					Namespace: "foo",
				},
				Spec: spec.APIProxySpec{
					Path: "/api/v1/accounts",
					Hosts: []spec.Host{
						{
							Name: "https://www.google.com",
						},
						{
							Name: "http://kubernetes.default.svc.cluster.local",
						},
					},
					Service: spec.Service{
						Name:      "bar",
						Namespace: "foo",
						Port:      8080,
					},
					Plugins: []spec.Plugin{
						{
							Name: "apikey",
						},
					},
				},
			},
		},
		{
			Source: &http.Request{
				URL: one,
			},
			Target: spec.APIProxy{
				TypeMeta: unversioned.TypeMeta{},
				ObjectMeta: api.ObjectMeta{
					Name:      "exampleAPIProxyOne",
					Namespace: "foo",
				},
				Spec: spec.APIProxySpec{
					Path:   "/api/v1/accounts",
					Target: "/foo/bar",
					Hosts: []spec.Host{
						{
							Name: "https://www.google.com",
						},
						{
							Name: "http://kubernetes.default.svc.cluster.local",
						},
					},
					Service: spec.Service{
						Name:      "bar",
						Namespace: "foo",
						Port:      8080,
					},
					Plugins: []spec.Plugin{
						{
							Name: "apikey",
						},
					},
					SSL: spec.SSL{
						SecretName: "mysecretname",
					},
				},
			},
		},
		{
			Source: &http.Request{
				URL: two,
			},
			Target: spec.APIProxy{
				TypeMeta: unversioned.TypeMeta{},
				ObjectMeta: api.ObjectMeta{
					Name:      "exampleAPIProxyOne",
					Namespace: "foo",
				},
				Spec: spec.APIProxySpec{
					Path:   "/api/v1/accounts",
					Target: "/foo/bar",
					Hosts: []spec.Host{
						{
							Name: "https://www.google.com",
						},
						{
							Name: "http://kubernetes.default.svc.cluster.local",
						},
					},
					Service: spec.Service{
						Name:      "my-service",
						Namespace: "foo",
						Port:      8080,
					},
					Plugins: []spec.Plugin{
						{
							Name: "apikey",
						},
					},
				},
			},
		},
		{
			Source: &http.Request{
				URL: three,
			},
			Target: spec.APIProxy{
				TypeMeta: unversioned.TypeMeta{},
				ObjectMeta: api.ObjectMeta{
					Name:      "exampleAPIProxyOne",
					Namespace: "foo",
				},
				Spec: spec.APIProxySpec{
					Path:   "/api/v1/accounts",
					Target: "/foo/bar",
					Service: spec.Service{
						Name:      "my-service",
						Namespace: "foo",
						Port:      8080,
					},
				},
			},
		},
	}

}
