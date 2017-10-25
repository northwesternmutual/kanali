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

package v2

// import (
// 	"testing"
//
// 	"github.com/stretchr/testify/assert"
// 	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
// )
//
// func TestGetKeyBindingStore(t *testing.T) {
// 	assert := assert.New(t)
// 	store := BindingStore
//
// 	store.Clear()
// 	assert.Equal(0, len(store.bindingMap), "store should be empty")
//
// 	v := BindingFactory{}
// 	var i interface{} = &v
// 	_, ok := i.(Store)
// 	assert.True(ok, "BindingFactory does not implement the Store interface")
// }
//
// func TestAPIKeyBindingSet(t *testing.T) {
// 	assert := assert.New(t)
// 	store := BindingStore
// 	keyBindingList := getTestAPIKeyBindingList()
//
// 	store.Clear()
// 	store.Set(keyBindingList.Items[0])
// 	err := store.Set(5)
// 	assert.Equal(err.Error(), "grrr - you're only allowed add api key bindings to the api key binding store.... duh", "wrong error")
// 	assert.Equal(keyBindingList.Items[0], store.bindingMap["foo"]["api-proxy-one"], "binding should exist")
// 	assert.Equal(1, len(store.bindingMap["foo"]["api-proxy-one"].Spec.Keys[0].SubpathTree.Children), "subpath incorrect")
// 	assert.Equal(1, len(store.bindingMap["foo"]["api-proxy-one"].Spec.Keys[0].SubpathTree.Children["foo"].Children), "subpath incorrect")
// 	store.Set(keyBindingList.Items[1])
// 	store.Set(keyBindingList.Items[2])
// 	assert.Equal(keyBindingList.Items[1], store.bindingMap["foo"]["api-proxy-two"], "binding should exist")
// 	assert.Equal(keyBindingList.Items[2], store.bindingMap["bar"]["api-proxy-three"], "binding should exist")
// 	store.Set(keyBindingList.Items[3])
// 	assert.Equal(keyBindingList.Items[3], store.bindingMap["foo"]["api-proxy-four"], "bidning should exist")
// }
//
// func TestAPIKeyBindingUpdate(t *testing.T) {
// 	assert := assert.New(t)
// 	store := BindingStore
// 	keyBindingList := getTestAPIKeyBindingList()
//
// 	store.Clear()
// 	store.Update(keyBindingList.Items[0], keyBindingList.Items[0])
// 	err := store.Update(5, 5)
// 	assert.Equal(err.Error(), "grrr - you're only allowed add api key bindings to the api key binding store.... duh", "wrong error")
// 	assert.Equal(keyBindingList.Items[0], store.bindingMap["foo"]["api-proxy-one"], "binding should exist")
// 	assert.Equal(1, len(store.bindingMap["foo"]["api-proxy-one"].Spec.Keys[0].SubpathTree.Children), "subpath incorrect")
// 	assert.Equal(1, len(store.bindingMap["foo"]["api-proxy-one"].Spec.Keys[0].SubpathTree.Children["foo"].Children), "subpath incorrect")
// 	store.Update(keyBindingList.Items[1], keyBindingList.Items[1])
// 	store.Update(keyBindingList.Items[2], keyBindingList.Items[2])
// 	assert.Equal(keyBindingList.Items[1], store.bindingMap["foo"]["api-proxy-two"], "binding should exist")
// 	assert.Equal(keyBindingList.Items[2], store.bindingMap["bar"]["api-proxy-three"], "binding should exist")
// 	store.Update(keyBindingList.Items[3], keyBindingList.Items[3])
// 	assert.Equal(keyBindingList.Items[3], store.bindingMap["foo"]["api-proxy-four"], "bidning should exist")
// }
//
// func TestAPIKeyBindingIsEmpty(t *testing.T) {
// 	assert := assert.New(t)
// 	store := BindingStore
// 	keyBindingList := getTestAPIKeyBindingList()
//
// 	store.Clear()
// 	assert.True(store.IsEmpty())
// 	store.Set(keyBindingList.Items[0])
// 	assert.False(store.IsEmpty())
// 	store.Clear()
// 	assert.True(store.IsEmpty())
//
// 	store.Set(keyBindingList.Items[0])
// 	assert.False(store.IsEmpty())
// 	store.Delete(keyBindingList.Items[0])
// 	assert.True(store.IsEmpty())
// }
//
// func TestAPIKeyBindingClear(t *testing.T) {
// 	assert := assert.New(t)
// 	store := BindingStore
// 	keyBindingList := getTestAPIKeyBindingList()
//
// 	store.Set(keyBindingList.Items[0])
// 	store.Clear()
// 	assert.Equal(0, len(store.bindingMap), "store should be empty")
// }
//
// func TestAPIKey(t *testing.T) {
// 	assert := assert.New(t)
// 	keyBindingList := getTestAPIKeyBindingList()
//
// 	assert.Equal(keyBindingList.Items[0].Spec.Keys[0], *keyBindingList.Items[0].GetAPIKey("franks-api-key"), "api keys should be equal")
// 	assert.Nil(keyBindingList.Items[0].GetAPIKey("bryans-api-key"), "apikey should be nil")
//
// 	assert.Equal(keyBindingList.Items[0].Spec.Keys[0], *keyBindingList.Items[0].GetAPIKey("franKs-aPi-Key"), "api keys should be equal")
// 	assert.Nil(keyBindingList.Items[0].GetAPIKey("bryans-api-key"), "apikey should be nil")
// }
//
// func TestGetRule(t *testing.T) {
// 	assert := assert.New(t)
// 	store := BindingStore
// 	keyBindingList := getTestAPIKeyBindingList()
// 	message := "rule note expected"
//
// 	store.Clear()
// 	store.Set(keyBindingList.Items[0])
// 	store.Set(keyBindingList.Items[1])
// 	store.Set(keyBindingList.Items[2])
// 	store.Set(keyBindingList.Items[3])
// 	store.Set(keyBindingList.Items[4])
// 	untypedResult, _ := store.Get("api-proxy-one", "foo")
// 	result, _ := untypedResult.(APIKeyBinding)
// 	assert.Equal(Rule{Global: true}, result.Spec.Keys[0].GetRule("/foo"), message)
// 	assert.Equal(Rule{Global: true}, result.Spec.Keys[0].GetRule("/foo/car"), message)
// 	assert.NotEqual(Rule{Global: true}, result.Spec.Keys[0].GetRule("/foo/bar"), message)
// 	assert.Equal(Rule{
// 		Granular: &GranularProxy{
// 			Verbs: []string{
// 				"POST",
// 				"GET",
// 			},
// 		},
// 	}, result.Spec.Keys[0].GetRule("/foo/bar"), message)
// 	assert.Equal(Rule{
// 		Granular: &GranularProxy{
// 			Verbs: []string{
// 				"POST",
// 				"GET",
// 			},
// 		},
// 	}, result.Spec.Keys[0].GetRule("/foo/bar/car"), message)
//
// 	untypedResult, _ = store.Get("api-proxy-two", "foo")
// 	result, _ = untypedResult.(APIKeyBinding)
// 	assert.Equal(Rule{Global: true}, result.Spec.Keys[0].GetRule("/foo"), message)
// 	assert.Equal(Rule{Global: true}, result.Spec.Keys[0].GetRule("/foo/bar"), message)
// 	assert.Equal(Rule{Global: false}, result.Spec.Keys[0].GetRule("/bar"), message)
//
// 	untypedResult, _ = store.Get("api-proxy-three", "bar")
// 	result, _ = untypedResult.(APIKeyBinding)
// 	assert.Equal(Rule{Global: true}, result.Spec.Keys[0].GetRule("/foo"), message)
// 	assert.Equal(Rule{Global: true}, result.Spec.Keys[0].GetRule("/foo/car"), message)
// 	assert.NotEqual(Rule{Global: true}, result.Spec.Keys[0].GetRule("/foo/bar"), message)
// 	assert.Equal(Rule{
// 		Granular: &GranularProxy{
// 			Verbs: []string{
// 				"POST",
// 				"GET",
// 			},
// 		},
// 	}, result.Spec.Keys[0].GetRule("/foo/bar"), message)
// 	assert.Equal(Rule{
// 		Granular: &GranularProxy{
// 			Verbs: []string{
// 				"POST",
// 				"GET",
// 			},
// 		},
// 	}, result.Spec.Keys[0].GetRule("/foo/bar/car"), message)
// 	assert.Equal(Rule{Global: true}, result.Spec.Keys[0].GetRule("/bar"), message)
//
// 	untypedResult, _ = store.Get("api-proxy-four", "foo")
// 	result, _ = untypedResult.(APIKeyBinding)
// 	assert.Equal(Rule{Global: true}, result.Spec.Keys[0].GetRule("/bar"), message)
// 	assert.Equal(Rule{Global: true}, result.Spec.Keys[0].GetRule("/foo"), message)
// 	assert.Equal(Rule{Global: true}, result.Spec.Keys[0].GetRule("/foo/bar"), message)
//
// 	untypedResult, _ = store.Get("api-proxy-five", "foo")
// 	result, _ = untypedResult.(APIKeyBinding)
// 	assert.Equal(Rule{Global: false}, result.Spec.Keys[0].GetRule("/bar"), message)
// 	assert.Equal(Rule{Global: false}, result.Spec.Keys[0].GetRule("/foo"), message)
// 	assert.Equal(Rule{Global: false}, result.Spec.Keys[0].GetRule("/foo/bar"), message)
//
// }
//
// func TestAPIKeyBindingGet(t *testing.T) {
// 	assert := assert.New(t)
// 	store := BindingStore
// 	keyList := getTestAPIKeyBindingList()
// 	message := "empty store should have no key bindings"
//
// 	store.Clear()
// 	store.Set(keyList.Items[0])
// 	store.Set(keyList.Items[1])
// 	store.Set(keyList.Items[2])
// 	_, err := store.Get(keyList.Items[2])
// 	assert.Equal(err.Error(), "should only pass the proxy name and namespace name", "wrong error")
// 	_, err = store.Get(5, "")
// 	assert.Equal(err.Error(), "proxy name should be a string", "wrong error")
// 	_, err = store.Get("", 5)
// 	assert.Equal(err.Error(), "namespace should be a string", "wrong error")
// 	result, _ := store.Get("", "")
// 	assert.Nil(result, message)
// 	result, _ = store.Get("api-proxy-one", "frank")
// 	assert.Nil(result, message)
// 	result, _ = store.Get("api-proxy-six", "foo")
// 	assert.Nil(result, message)
// 	result, _ = store.Get("api-proxy-one", "foo")
// 	assert.Equal(keyList.Items[0], result, message)
// 	result, _ = store.Get("api-proxy-two", "foo")
// 	assert.Equal(keyList.Items[1], result, message)
// 	result, _ = store.Get("api-proxy-three", "bar")
// 	assert.Equal(keyList.Items[2], result, message)
// }
//
// func TestAPIKeyBindingDelete(t *testing.T) {
// 	assert := assert.New(t)
// 	store := BindingStore
// 	keyList := getTestAPIKeyBindingList()
// 	message := "empty store should have no keys"
//
// 	_, err := store.Delete(5)
// 	assert.Equal(err.Error(), "there's no way this api key binding could've gotten in here", "wrong error")
// 	store.Clear()
// 	store.Set(keyList.Items[0])
// 	result, _ := store.Delete(keyList.Items[2])
// 	assert.Nil(result, message)
// 	store.Set(keyList.Items[3])
// 	store.Set(keyList.Items[1])
// 	result, _ = store.Delete(keyList.Items[3])
// 	assert.Equal(keyList.Items[3], result, message)
// 	result, _ = store.Delete(nil)
// 	assert.Nil(result, message)
// 	result, _ = store.Delete(keyList.Items[2])
// 	assert.Nil(result, message)
// 	result, _ = store.Delete(keyList.Items[1])
// 	assert.Equal(keyList.Items[1], result, message)
// 	result, _ = store.Delete(keyList.Items[1])
// 	assert.Nil(result, message)
// 	result, _ = store.Get("api-proxy-two", "foo")
// 	assert.Nil(result, message)
// 	result, _ = store.Get("api-proxy-four", "foo")
// 	assert.Nil(result, message)
// 	assert.Equal(1, len(store.bindingMap), message)
// 	store.Set(keyList.Items[1])
// 	result, _ = store.Get("api-proxy-two", "foo")
// 	assert.Equal(keyList.Items[1], result, message)
// 	assert.Equal(1, len(store.bindingMap), message)
// 	store.Set(keyList.Items[2])
// 	result, _ = store.Delete(keyList.Items[2])
// 	assert.Equal(keyList.Items[2], result, message)
// 	result, _ = store.Get("api-proxy-two", "foo")
// 	assert.NotNil(result, message)
// 	result, _ = store.Get("api-proxy-three", "bar")
// 	assert.Nil(result, message)
// }
//
// func getTestAPIKeyBindingList() *APIKeyBindingList {
//
// 	return &APIKeyBindingList{
// 		ListMeta: metav1.ListMeta{},
// 		Items: []APIKeyBinding{
// 			{
// 				ObjectMeta: metav1.ObjectMeta{
// 					Name:      "abc123",
// 					Namespace: "foo",
// 				},
// 				Spec: APIKeyBindingSpec{
// 					APIProxyName: "api-proxy-one",
// 					Keys: []Key{
// 						{
// 							Name:        "franks-api-key",
// 							DefaultRule: Rule{},
// 							Subpaths: []*Path{
// 								{
// 									Path: "/foo",
// 									Rule: Rule{
// 										Global: true,
// 									},
// 								},
// 								{
// 									Path: "foo/bar",
// 									Rule: Rule{
// 										Granular: &GranularProxy{
// 											Verbs: []string{
// 												"POST",
// 												"GET",
// 											},
// 										},
// 									},
// 								},
// 							},
// 						},
// 					},
// 				},
// 			},
// 			{
// 				ObjectMeta: metav1.ObjectMeta{
// 					Name:      "abc123",
// 					Namespace: "foo",
// 				},
// 				Spec: APIKeyBindingSpec{
// 					APIProxyName: "api-proxy-two",
// 					Keys: []Key{
// 						{
// 							Name:        "franks-api-key",
// 							DefaultRule: Rule{},
// 							Subpaths: []*Path{
// 								{
// 									Path: "/foo",
// 									Rule: Rule{
// 										Global: true,
// 									},
// 								},
// 							},
// 						},
// 					},
// 				},
// 			},
// 			{
// 				ObjectMeta: metav1.ObjectMeta{
// 					Name:      "abc123",
// 					Namespace: "bar",
// 				},
// 				Spec: APIKeyBindingSpec{
// 					APIProxyName: "api-proxy-three",
// 					Keys: []Key{
// 						{
// 							Name: "franks-api-key",
// 							DefaultRule: Rule{
// 								Global: true,
// 							},
// 							Subpaths: []*Path{
// 								{
// 									Path: "/foo",
// 									Rule: Rule{
// 										Global: true,
// 									},
// 								},
// 								{
// 									Path: "/foo/bar",
// 									Rule: Rule{
// 										Granular: &GranularProxy{
// 											Verbs: []string{
// 												"POST",
// 												"GET",
// 											},
// 										},
// 									},
// 								},
// 							},
// 						},
// 					},
// 				},
// 			},
// 			{
// 				ObjectMeta: metav1.ObjectMeta{
// 					Name:      "abc123",
// 					Namespace: "foo",
// 				},
// 				Spec: APIKeyBindingSpec{
// 					APIProxyName: "api-proxy-four",
// 					Keys: []Key{
// 						{
// 							Name: "franks-api-key",
// 							DefaultRule: Rule{
// 								Global: true,
// 							},
// 						},
// 					},
// 				},
// 			},
// 			{
// 				ObjectMeta: metav1.ObjectMeta{
// 					Name:      "abc123",
// 					Namespace: "foo",
// 				},
// 				Spec: APIKeyBindingSpec{
// 					APIProxyName: "api-proxy-five",
// 					Keys: []Key{
// 						{
// 							Name: "franks-api-key",
// 						},
// 					},
// 				},
// 			},
// 		},
// 	}
//
// }
