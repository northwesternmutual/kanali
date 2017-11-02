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

import (
	"testing"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
  "github.com/northwesternmutual/kanali/pkg/apis/kanali.io/v2"
)


func TestApiKeyBindingStore(t *testing.T) {
	_, ok := ApiKeyBindingStore().(ApiKeyBindingStoreInterface)
	assert.True(t, ok)
}

func TestApiKeyBindingStoreSet(t *testing.T) {
	defer ApiKeyBindingStore().Clear()

  testSubpathOne := getTestSubpaths("/foo", true)
  testSubpathTwo := getTestSubpaths("/foo/bar", true, "GET", "DELETE")

  testKeyOne := getTestKey("example-one", v2.Rule{Granular: v2.GranularProxy{Verbs: []string{"POST", "PUT"}}},
    testSubpathOne,
    testSubpathTwo,
  )

  testApiKeyBindingOne := getTestApiKeyBinding("example-one", "foo", testKeyOne)

  ApiKeyBindingStore().Set(testApiKeyBindingOne)

  assert.Equal(t, 1, len(apiKeyBindingStore.apiKeyBindingMap["foo"]))
  assert.Equal(t, 1, len(apiKeyBindingStore.apiKeyBindingMap["foo"]["example-one"]))
  assert.Equal(t, testKeyOne, apiKeyBindingStore.apiKeyBindingMap["foo"]["example-one"]["example-one"].key)
  assert.Equal(t, 1, len(apiKeyBindingStore.apiKeyBindingMap["foo"]["example-one"]["example-one"].subpathTree.children))
  assert.Equal(t, 1, len(apiKeyBindingStore.apiKeyBindingMap["foo"]["example-one"]["example-one"].subpathTree.children["foo"].children))
  assert.Equal(t, &testSubpathOne, apiKeyBindingStore.apiKeyBindingMap["foo"]["example-one"]["example-one"].subpathTree.children["foo"].value)
  assert.Equal(t, &testSubpathTwo, apiKeyBindingStore.apiKeyBindingMap["foo"]["example-one"]["example-one"].subpathTree.children["foo"].children["bar"].value)
}

func TestApiKeyBindingStoreClear(t *testing.T) {
  defer ApiKeyBindingStore().Clear()

  assert.Equal(t, 0, len(apiKeyBindingStore.apiKeyBindingMap))

  testSubpathOne := getTestSubpaths("/foo", true)
  testKeyOne := getTestKey("example-one", v2.Rule{Granular: v2.GranularProxy{Verbs: []string{"POST", "PUT"}}}, testSubpathOne)
  testApiKeyBindingOne := getTestApiKeyBinding("example-one", "foo", testKeyOne)

  ApiKeyBindingStore().Set(testApiKeyBindingOne)

  assert.Equal(t, 1, len(apiKeyBindingStore.apiKeyBindingMap))
  ApiKeyBindingStore().Clear()
  assert.Equal(t, 0, len(apiKeyBindingStore.apiKeyBindingMap))
}


func TestAPIKeyBindingUpdate(t *testing.T) {
  defer ApiKeyBindingStore().Clear()

  testSubpathOne := getTestSubpaths("/foo", true)
  testSubpathTwo := getTestSubpaths("/foo/bar", true, "GET", "DELETE")

  testKeyOneOld := getTestKey("example-one", v2.Rule{Granular: v2.GranularProxy{Verbs: []string{"POST", "PUT"}}},
    testSubpathOne,
    testSubpathTwo,
  )
  testKeyOneNew := getTestKey("example-one", v2.Rule{Granular: v2.GranularProxy{Verbs: []string{"POST", "PUT"}}},
    testSubpathTwo,
  )

  testApiKeyBindingOneOld := getTestApiKeyBinding("example-one", "foo", testKeyOne)
  testApiKeyBindingOneNew := getTestApiKeyBinding("example-one", "foo", testKeyOneNew)

  ApiKeyBindingStore().Set(testApiKeyBindingOneOld)
  ApiKeyBindingStore().Update(testApiKeyBindingOneOld, testApiKeyBindingOneNew)
}

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

func getTestKey(name string, rule v2.Rule, subpaths ...v2.Path) v2.Key {
  return v2.Key{
    Name: name,
    DefaultRule: rule,
    Subpaths: subpaths,
  }
}

func getTestSubpaths(path string, global bool, methods ...string) v2.Path {
  var rule v2.Rule

  if global {
    rule.Global = true
  } else {
    rule.Granular = v2.GranularProxy{
      Verbs: methods,
    }
  }

  return v2.Path{
    Path: path,
    Rule: rule,
  }
}

func getTestApiKeyBinding(name, namespace string, keys ...v2.Key) *v2.ApiKeyBinding {
	return &v2.ApiKeyBinding{
    ObjectMeta: metav1.ObjectMeta{
      Name:      name,
      Namespace: namespace,
    },
    Spec: v2.ApiKeyBindingSpec{
      Keys: keys,
    },
  }
}
