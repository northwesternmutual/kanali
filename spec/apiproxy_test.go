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

package spec

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/unversioned"
)

func TestAPIProxyGetProxyStore(t *testing.T) {
	assert := assert.New(t)
	store := ProxyStore

	store.Clear()
	assert.Equal(0, len(store.proxyTree.Children), "store should be empty")
	assert.Nil(store.proxyTree.Value, "empty store should have no value")

	v := ProxyFactory{}
	var i interface{} = &v
	_, ok := i.(Store)
	assert.True(ok, "ProxyFactory does not implement the Store interface")

}

func TestAPIProxySet(t *testing.T) {
	assert := assert.New(t)
	store := ProxyStore
	proxyList := getTestAPIProxyList()

	store.Clear()
	store.Set(proxyList.Proxies[0])
	store.Set(proxyList.Proxies[1])
	store.Set(proxyList.Proxies[2])
	err := store.Set(APIKey{})
	assert.Equal(err.Error(), "grrr - you're only allowed add api proxies to the api proxy store.... duh", "wrong error")
	assert.Equal(proxyList.Proxies[0], *store.proxyTree.Children["api"].Children["v1"].Children["accounts"].Value, "proxy should exist")
	assert.Equal(proxyList.Proxies[1], *store.proxyTree.Children["api"].Children["v1"].Children["field"].Value, "proxy should exist")
	assert.Equal(proxyList.Proxies[2], *store.proxyTree.Children["api"].Value, "proxy should exist")
}

func TestAPIProxyClear(t *testing.T) {
	assert := assert.New(t)
	store := ProxyStore
	proxyList := getTestAPIProxyList()

	store.Set(proxyList.Proxies[0])
	store.Clear()
	assert.Equal(0, len(store.proxyTree.Children), "store should be empty")
}

func TestDeletePreviousProxy(t *testing.T) {
	ProxyStore.Clear()
	defer ProxyStore.Clear()
	proxyList := getTestAPIProxyList()
	ProxyStore.Set(proxyList.Proxies[0])
	ProxyStore.Set(proxyList.Proxies[1])
	ProxyStore.Set(proxyList.Proxies[2])
	ProxyStore.Set(proxyList.Proxies[3])

	assert.Equal(t, len(ProxyStore.proxyTree.Children["api"].Children["v1"].Children), 2)
	ProxyStore.proxyTree.deletePreviousProxy(proxyList.Proxies[1])
	assert.False(t, ProxyStore.IsEmpty())
	assert.Equal(t, len(ProxyStore.proxyTree.Children), 1)
	assert.Equal(t, len(ProxyStore.proxyTree.Children["api"].Children), 1)
	assert.Equal(t, len(ProxyStore.proxyTree.Children["api"].Children["v1"].Children), 1)
	untyped, _ := ProxyStore.Get("/api/v1/field")
	typed, _ := untyped.(APIProxy)
	assert.Equal(t, typed.Spec.Path, "/api/v1")
	untyped, _ = ProxyStore.Get("/api/v1")
	typed, _ = untyped.(APIProxy)
	assert.Equal(t, typed.Spec.Path, "/api/v1")
	ProxyStore.proxyTree.deletePreviousProxy(proxyList.Proxies[3])
	untyped, _ = ProxyStore.Get("/api/v1")
	typed, _ = untyped.(APIProxy)
	untyped, _ = ProxyStore.Get("/api/v1")
	typed, _ = untyped.(APIProxy)
	assert.Equal(t, typed.Spec.Path, "/api")
	assert.Equal(t, len(ProxyStore.proxyTree.Children), 1)
	assert.Equal(t, len(ProxyStore.proxyTree.Children["api"].Children), 1)
	assert.Equal(t, len(ProxyStore.proxyTree.Children["api"].Children["v1"].Children), 1)
	untyped, _ = ProxyStore.Get("/api/v1/accounts")
	assert.NotNil(t, untyped)

	ProxyStore.proxyTree.deletePreviousProxy(proxyList.Proxies[0])
	assert.False(t, ProxyStore.IsEmpty())
	ProxyStore.proxyTree.deletePreviousProxy(proxyList.Proxies[2])
	assert.True(t, ProxyStore.IsEmpty())
}

func TestAPIProxyIsEmpty(t *testing.T) {
	assert := assert.New(t)
	store := ProxyStore
	proxyList := getTestAPIProxyList()

	store.Clear()
	assert.True(store.IsEmpty())
	store.Set(proxyList.Proxies[0])
	assert.False(store.IsEmpty())
	store.Clear()
	assert.True(store.IsEmpty())
}

func TestAPIProxyGet(t *testing.T) {
	assert := assert.New(t)
	store := ProxyStore
	proxyList := getTestAPIProxyList()

	store.Clear()
	store.Set(proxyList.Proxies[0])
	store.Set(proxyList.Proxies[1])
	store.Set(proxyList.Proxies[2])
	_, err := store.Get("", "")
	assert.Equal(err.Error(), "should only pass the path of the proxy", "wrong error")
	_, err = store.Get(5)
	assert.Equal(err.Error(), "when retrieving a proxy, use the proxy path", "wrong error")
	result, _ := store.Get("")
	assert.Nil(result, "proxy should not be returned")
	result, _ = store.Get("foo")
	assert.Nil(result, "proxy should not be returned")
	result, _ = store.Get("bar")
	assert.Nil(result, "proxy should not be returned")
	// result, _ = store.Get("api/v1")
	// assert.Nil(result, "proxy should not be returned")
	result, _ = store.Get("api/v1")
	assert.Equal(proxyList.Proxies[2], result, "proxy should be returned")

	result, _ = store.Get("foo/bar")
	assert.Nil(result, "proxy should not be returned")
	result, _ = store.Get("api/v1/accounts")
	assert.Equal(proxyList.Proxies[0], result, "proxy should be returned")
	result, _ = store.Get("api/v1/accounts/foo/bar")
	assert.Equal(proxyList.Proxies[0], result, "proxy should be returned")
	result, _ = store.Get("/api/v1/field")
	assert.Equal(proxyList.Proxies[1], result, "proxy should be returned")
	result, _ = store.Get("api/v1/field/foo/bar")
	assert.Equal(proxyList.Proxies[1], result, "proxy should be returned")
	result, _ = store.Get("api")
	assert.Equal(proxyList.Proxies[2], result, "proxy should be returned")
	result, _ = store.Get("api/foo")
	assert.Equal(proxyList.Proxies[2], result, "proxy should be returned")
}

func TestAPIProxyDelete(t *testing.T) {
	assert := assert.New(t)
	store := ProxyStore
	proxyList := getTestAPIProxyList()
	message := "proxy not deleted correctly"

	store.Clear()
	store.Set(proxyList.Proxies[0])
	store.Set(proxyList.Proxies[1])
	_, err := store.Delete(5)
	assert.Equal(err.Error(), "there's no way this api proxy could've gotten in here", "wrong error")
	result, _ := store.Delete(nil)
	assert.Nil(result, message)
	result, _ = store.Delete(proxyList.Proxies[2])
	assert.Nil(result, message)
	result, _ = store.Delete(proxyList.Proxies[3])
	assert.Nil(result, message)
	result, _ = store.Get("api/v1/field")
	assert.Equal(proxyList.Proxies[1], result, message)
	result, _ = store.Delete(proxyList.Proxies[1])
	assert.Equal(proxyList.Proxies[1], result, message)
	result, _ = store.Get("api/v1/field")
	assert.Nil(result, message)
	result, _ = store.Get("api/v1/field/foo")
	assert.Nil(result, message)
	result, _ = store.Get("api/v1/accounts")
	assert.NotNil(result, message)
	assert.Equal(1, len(store.proxyTree.Children["api"].Children["v1"].Children), message)
	store.Set(proxyList.Proxies[1])
	result, _ = store.Get("api/v1/field")
	assert.Equal(proxyList.Proxies[1], result, message)
	assert.Equal(2, len(store.proxyTree.Children["api"].Children["v1"].Children), message)
	store.Set(proxyList.Proxies[2])
	result, _ = store.Delete(proxyList.Proxies[2])
	assert.Equal(proxyList.Proxies[2], result, message)
	result, _ = store.Get("api/v1/accounts")
	assert.NotNil(result, message)
	result, _ = store.Get("api")
	assert.Nil(result, message)
}

func TestGetFileName(t *testing.T) {
	assert := assert.New(t)
	store := ProxyStore
	proxyList := getTestAPIProxyList()
	message := "file name not what expected"

	store.Clear()
	store.Set(proxyList.Proxies[0])
	store.Set(proxyList.Proxies[1])
	store.Set(proxyList.Proxies[2])
	store.Set(proxyList.Proxies[2])
	assert.Equal("jwt", proxyList.Proxies[0].Spec.Plugins[1].GetFileName(), message)
	assert.Equal("apikey_1.0.0", proxyList.Proxies[0].Spec.Plugins[0].GetFileName(), message)
}

func TestGetSSLCertificates(t *testing.T) {
	assert := assert.New(t)
	ProxyFactory := ProxyStore
	proxyList := getTestAPIProxyList()
	message := "ssl object not as expected"

	ProxyFactory.Clear()
	ProxyFactory.Set(proxyList.Proxies[0])
	ProxyFactory.Set(proxyList.Proxies[1])
	ProxyFactory.Set(proxyList.Proxies[2])
	ProxyFactory.Set(proxyList.Proxies[3])

	untypedResult, _ := ProxyFactory.Get("/api/v1")
	result, _ := untypedResult.(APIProxy)
	assert.Equal(SSL{}, *result.GetSSLCertificates(""), message)
	untypedResult, _ = ProxyFactory.Get("api/v1/accounts")
	result, _ = untypedResult.(APIProxy)
	assert.Equal(SSL{"mySecret"}, *result.GetSSLCertificates(""), message)
	untypedResult, _ = ProxyFactory.Get("api/v1/accounts")
	result, _ = untypedResult.(APIProxy)
	assert.Equal(SSL{"mySecretTwo"}, *result.GetSSLCertificates("foo.bar.com"), message)
	untypedResult, _ = ProxyFactory.Get("api/v1/accounts")
	result, _ = untypedResult.(APIProxy)
	assert.Equal(SSL{"mySecret"}, *result.GetSSLCertificates("bar.foo.com"), message)
}

func TestNormalize(t *testing.T) {
	p1 := APIProxy{
		Spec: APIProxySpec{
			Path:   "///foo/bar///car",
			Target: "foo///bar//car",
		},
	}
	p2 := APIProxy{
		Spec: APIProxySpec{
			Path:   "",
			Target: "///",
		},
	}
	normalize(&p1)
	normalize(&p2)

	assert.Equal(t, p1.Spec.Path, "/foo/bar/car")
	assert.Equal(t, p1.Spec.Target, "/foo/bar/car")
	assert.Equal(t, p2.Spec.Path, "/")
	assert.Equal(t, p2.Spec.Target, "/")
}

func getTestAPIProxyList() *APIProxyList {

	return &APIProxyList{
		TypeMeta: unversioned.TypeMeta{},
		ListMeta: unversioned.ListMeta{},
		Proxies: []APIProxy{
			{
				TypeMeta: unversioned.TypeMeta{},
				ObjectMeta: api.ObjectMeta{
					Name:      "exampleAPIProxyOne",
					Namespace: "foo",
				},
				Spec: APIProxySpec{
					Path:   "/api/v1/accounts",
					Target: "/",
					Hosts: []Host{
						{
							Name: "foo.bar.com",
							SSL: SSL{
								SecretName: "mySecretTwo",
							},
						},
						{
							Name: "bar.foo.com",
						},
					},
					Service: Service{
						Name:      "my-service",
						Namespace: "foo",
						Port:      8080,
					},
					Plugins: []Plugin{
						{
							Name:    "apikey",
							Version: "1.0.0",
						},
						{
							Name: "jwt",
						},
					},
					SSL: SSL{
						SecretName: "mySecret",
					},
				},
			},
			{
				TypeMeta: unversioned.TypeMeta{},
				ObjectMeta: api.ObjectMeta{
					Name:      "exampleAPIProxyTwo",
					Namespace: "foo",
				},
				Spec: APIProxySpec{
					Path:   "/api/v1/field",
					Target: "/",
					Hosts: []Host{
						{
							Name: "https://www.google.com",
						},
						{
							Name: "http://kubernetes.default.svc.cluster.local",
						},
					},
					Service: Service{
						Name:      "my-service",
						Namespace: "foo",
						Port:      8080,
					},
					Plugins: []Plugin{
						{
							Name: "apikey",
						},
					},
					SSL: SSL{
						SecretName: "mySecret",
					},
				},
			},
			{
				TypeMeta: unversioned.TypeMeta{},
				ObjectMeta: api.ObjectMeta{
					Name:      "exampleAPIProxyThree",
					Namespace: "foo",
				},
				Spec: APIProxySpec{
					Path:   "/api",
					Target: "/",
					Hosts: []Host{
						{
							Name: "https://www.google.com",
						},
						{
							Name: "http://kubernetes.default.svc.cluster.local",
						},
					},
					Service: Service{
						Name:      "my-service",
						Namespace: "foo",
						Port:      8080,
					},
					Plugins: []Plugin{
						{
							Name: "apikey",
						},
						{
							Name: "jwt",
						},
						{
							Name: "quota",
						},
					},
					SSL: SSL{
						SecretName: "mySecret",
					},
				},
			},
			{
				TypeMeta: unversioned.TypeMeta{},
				ObjectMeta: api.ObjectMeta{
					Name:      "exampleAPIProxyFour",
					Namespace: "foo",
				},
				Spec: APIProxySpec{
					Path:   "/api/v1",
					Target: "/",
					Hosts: []Host{
						{
							Name: "https://www.google.com",
						},
						{
							Name: "http://kubernetes.default.svc.cluster.local",
						},
					},
					Service: Service{
						Name:      "my-service",
						Namespace: "foo",
						Port:      8080,
					},
					Plugins: []Plugin{
						{
							Name: "apikey",
						},
						{
							Name: "jwt",
						},
						{
							Name: "quota",
						},
					},
				},
			},
		},
	}

}
