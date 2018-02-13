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

package apikey

// import (
//   "os"
//   "fmt"
//   "testing"
// )

// func TestOnRequest(t *testing.T) {
// 	assert := assert.New(t)
//
// 	assert.Nil(Plugin.OnRequest(context.Background(), &metrics.Metrics{}, getTestAPIProxy(), &http.Request{
// 		Method: "OPTIONS",
// 	}, opentracing.StartSpan("test span")))
// 	assert.Nil(Plugin.OnRequest(context.Background(), &metrics.Metrics{}, getTestAPIProxy(), &http.Request{
// 		Method: "options",
// 	}, opentracing.StartSpan("test span")))
//
// 	assert.Equal("apikey not found in request", Plugin.OnRequest(context.Background(), &metrics.Metrics{}, spec.APIProxy{}, &http.Request{}, opentracing.StartSpan("test span")).Error(), "should have thrown error")
//
// 	viper.SetDefault(flagPluginsAPIKeyHeaderKey.GetLong(), "apikey")
//
// 	u, _ := url.Parse("http://host.com/api/v1/accounts")
//
// 	assert.Equal("apikey not found in k8s cluster", Plugin.OnRequest(context.Background(), &metrics.Metrics{}, getTestAPIProxy(), &http.Request{
// 		Header: http.Header{
// 			"Apikey": []string{"myapikey"},
// 		},
// 		URL: u,
// 	}, opentracing.StartSpan("test span")).Error(), "should have thrown error")
//
// 	apikeyStore := spec.KeyStore
// 	apikeyStore.Set(getTestAPIKey())
//
// 	assert.Equal("no binding found for associated APIProxy", Plugin.OnRequest(context.Background(), &metrics.Metrics{}, getTestAPIProxy(), &http.Request{
// 		Header: http.Header{
// 			"Apikey": []string{"myapikey"},
// 		},
// 		URL: u,
// 	}, opentracing.StartSpan("test span")).Error(), "should have thrown error")
//
// 	binding := getTestAPIKeyBinding()
// 	binding.Spec.Keys = []spec.Key{}
// 	apikeybindingStore := spec.BindingStore
// 	apikeybindingStore.Set(binding)
//
// 	assert.Equal("api key not authorized for this proxy", Plugin.OnRequest(context.Background(), &metrics.Metrics{}, getTestAPIProxy(), &http.Request{
// 		Header: http.Header{
// 			"Apikey": []string{"myapikey"},
// 		},
// 		URL: u,
// 	}, opentracing.StartSpan("test span")).Error(), "should have thrown error")
//
// 	spec.KanaliEndpoints = &api.Endpoints{
// 		TypeMeta:   unversioned.TypeMeta{},
// 		ObjectMeta: api.ObjectMeta{},
// 		Subsets: []api.EndpointSubset{
// 			{
// 				Addresses: []api.EndpointAddress{
// 					{
// 						IP: "1.2.3.4",
// 					},
// 				},
// 			},
// 		},
// 	}
//
// 	apikeybindingStore.Set(getTestAPIKeyBinding())
// 	assert.Nil(Plugin.OnRequest(context.Background(), &metrics.Metrics{}, getTestAPIProxy(), &http.Request{
// 		Header: http.Header{
// 			"Apikey": []string{"myapikey"},
// 		},
// 		URL: u,
// 	}, opentracing.StartSpan("test span")), "apikey should be authorized")
// }
//
// func TestOnResponse(t *testing.T) {
// 	assert := assert.New(t)
// 	assert.Nil(Plugin.OnResponse(context.Background(), &metrics.Metrics{}, spec.APIProxy{}, &http.Request{}, nil, opentracing.StartSpan("test span")))
// }
//
// func TestValidateAPIKey(t *testing.T) {
// 	assert := assert.New(t)
//
// 	assert.True(validateAPIKey(spec.Rule{
// 		Global: true,
// 		Granular: &spec.GranularProxy{
// 			Verbs: []string{},
// 		},
// 	}, "GET"), "rule should be authorized")
//
// 	assert.True(validateAPIKey(spec.Rule{
// 		Global: true,
// 		Granular: &spec.GranularProxy{
// 			Verbs: []string{
// 				"GET",
// 			},
// 		},
// 	}, "GET"), "rule should be authorized")
//
// 	assert.True(validateAPIKey(spec.Rule{
// 		Global: false,
// 		Granular: &spec.GranularProxy{
// 			Verbs: []string{
// 				"GET",
// 			},
// 		},
// 	}, "GET"), "rule should be authorized")
//
// 	assert.True(validateAPIKey(spec.Rule{
// 		Global: true,
// 		Granular: &spec.GranularProxy{
// 			Verbs: []string{
// 				"POST",
// 			},
// 		},
// 	}, "GET"), "rule should be authorized")
//
// 	assert.False(validateAPIKey(spec.Rule{
// 		Global: false,
// 		Granular: &spec.GranularProxy{
// 			Verbs: []string{
// 				"POST",
// 			},
// 		},
// 	}, "GET"), "rule should not be authorized")
// }
//
// func TestValidateGranularRules(t *testing.T) {
// 	assert := assert.New(t)
//
// 	assert.True(validateGranularRules("GET", &spec.GranularProxy{
// 		Verbs: []string{
// 			"GET",
// 		},
// 	}), "http method should be authorized")
//
// 	assert.True(validateGranularRules("get", &spec.GranularProxy{
// 		Verbs: []string{
// 			"GET",
// 			"POST",
// 			"PUT",
// 		},
// 	}), "http method should be authorized")
//
// 	assert.True(validateGranularRules("put", &spec.GranularProxy{
// 		Verbs: []string{
// 			"GET",
// 			"POST",
// 			"PUT",
// 		},
// 	}), "http method should be authorized")
//
// 	assert.False(validateGranularRules("GET", &spec.GranularProxy{
// 		Verbs: []string{
// 			"POST",
// 		},
// 	}), "http method should be authorized")
//
// 	assert.False(validateGranularRules("get", &spec.GranularProxy{
// 		Verbs: []string{
// 			"POST",
// 			"PUT",
// 		},
// 	}), "http method should be authorized")
//
// 	assert.False(validateGranularRules("", &spec.GranularProxy{
// 		Verbs: []string{
// 			"POST",
// 		},
// 	}), "http method should be authorized")
//
// 	assert.False(validateGranularRules("HTTP", &spec.GranularProxy{
// 		Verbs: []string{
// 			"POST",
// 		},
// 	}), "http method should be authorized")
// }
//
// func getTestAPIProxy() spec.APIProxy {
//
// 	return spec.APIProxy{
// 		TypeMeta: unversioned.TypeMeta{},
// 		ObjectMeta: api.ObjectMeta{
// 			Name:      "APIProxyone",
// 			Namespace: "foo",
// 		},
// 		Spec: spec.APIProxySpec{
// 			Path:   "/api/v1/accounts",
// 			Target: "/",
// 			Service: spec.Service{
// 				Name:      "my-service",
// 				Namespace: "foo",
// 				Port:      8080,
// 			},
// 			Plugins: []spec.Plugin{
// 				{
// 					Name: "apikey",
// 				},
// 			},
// 		},
// 	}
//
// }
//
// func getTestAPIKey() spec.APIKey {
//
// 	return spec.APIKey{
// 		TypeMeta: unversioned.TypeMeta{},
// 		ObjectMeta: api.ObjectMeta{
// 			Name:      "apikeyone",
// 			Namespace: "foo",
// 		},
// 		Spec: spec.APIKeySpec{
// 			APIKeyData: "myapikey",
// 		},
// 	}
//
// }
//
// func getTestAPIKeyBinding() spec.APIKeyBinding {
//
// 	return spec.APIKeyBinding{
// 		TypeMeta: unversioned.TypeMeta{},
// 		ObjectMeta: api.ObjectMeta{
// 			Name:      "apikeybindingone",
// 			Namespace: "foo",
// 		},
// 		Spec: spec.APIKeyBindingSpec{
// 			APIProxyName: "APIProxyone",
// 			Keys: []spec.Key{
// 				{
// 					Name: "apikeyone",
// 					DefaultRule: spec.Rule{
// 						Global: true,
// 					},
// 				},
// 			},
// 		},
// 	}
//
// }
