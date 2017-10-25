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

func TestApiProxyStore(t *testing.T) {
	_, ok := ApiProxyStore().(ApiProxyStoreInterface)
	assert.True(t, ok)
}

func TestApiProxyStoreSet(t *testing.T) {
  defer ApiProxyStore().Clear()

  testApiProxyOne := getTestApiProxy("example-one", "foo", "/api/v1", "/")
  testApiProxyTwo := getTestApiProxy("example-two", "foo", "/api/v2", "/")
  testApiProxyThree := getTestApiProxy("example-three", "foo", "/api", "/")
  testApiProxyFour := getTestApiProxy("example-four", "foo", "/apis", "/")
  testApiProxyFive := getTestApiProxy("example-five", "foo", "/", "/")

	assert.Nil(t, ApiProxyStore().Set(testApiProxyOne))
  assert.Nil(t, ApiProxyStore().Set(testApiProxyTwo))
  assert.Nil(t, ApiProxyStore().Set(testApiProxyThree))
  assert.Nil(t, ApiProxyStore().Set(testApiProxyFour))
  assert.Nil(t, ApiProxyStore().Set(testApiProxyFive))

	assert.Equal(t, testApiProxyOne, apiProxyStore.proxyTree.children["api"].children["v1"].value)
	assert.Equal(t, testApiProxyTwo, apiProxyStore.proxyTree.children["api"].children["v2"].value)
	assert.Equal(t, testApiProxyThree, apiProxyStore.proxyTree.children["api"].value)
  assert.Equal(t, testApiProxyFour, apiProxyStore.proxyTree.children["apis"].value)
  assert.Equal(t, testApiProxyFive, apiProxyStore.proxyTree.value)
}

func TestApiProxyStoreClear(t *testing.T) {
  assert.Nil(t, ApiProxyStore().Set(getTestApiProxy("example-one", "foo", "/api/v1", "/")))
  assert.Equal(t, 1, len(apiProxyStore.proxyTree.children))
	ApiProxyStore().Clear()
	assert.Equal(t, 0, len(apiProxyStore.proxyTree.children))
}

func TestApiProxyStoreUpdate(t *testing.T) {
	defer ApiProxyStore().Clear()

  testApiProxyOneOld := getTestApiProxy("example-one", "foo", "/api/v1", "/foo")
  testApiProxyTwoOld := getTestApiProxy("example-two", "foo", "/apis", "/")

  assert.Nil(t, ApiProxyStore().Set(testApiProxyOneOld))
  assert.Nil(t, ApiProxyStore().Set(testApiProxyTwoOld))

  testApiProxyOneNew := getTestApiProxy("example-one", "foo", "/api/v1", "/")
  testApiProxyTwoNew := getTestApiProxy("example-two", "foo", "/", "/")

	assert.Nil(t, ApiProxyStore().Update(testApiProxyOneOld, testApiProxyOneNew))
  assert.Nil(t, ApiProxyStore().Update(testApiProxyTwoOld, testApiProxyTwoNew))

	assert.Equal(t, testApiProxyOneNew, apiProxyStore.proxyTree.children["api"].children["v1"].value)
  assert.Equal(t, testApiProxyTwoNew, apiProxyStore.proxyTree.children[""].value)
  assert.Nil(t, apiProxyStore.proxyTree.children["apis"])

  ApiProxyStore().Clear()
  assert.Nil(t, ApiProxyStore().Set(testApiProxyOneOld))
  assert.Equal(t, ApiProxyStore().Update(testApiProxyOneOld, getTestApiProxy("example-three", "foo", "/api/v1", "/")).Error(), "ApiProxy with different ObjectMeta exists at this path")
}

func TestApiProxyStoreIsEmpty(t *testing.T) {
  defer ApiProxyStore().Clear()

  assert.True(t, ApiProxyStore().IsEmpty())
  assert.Nil(t, ApiProxyStore().Set(getTestApiProxy("example-one", "foo", "/api/v1", "/")))
  assert.False(t, ApiProxyStore().IsEmpty())
}

func TestApiProxyStoreGet(t *testing.T) {
  defer ApiProxyStore().Clear()

  testApiProxyOne := getTestApiProxy("example-one", "foo", "/api/v1", "/")
  testApiProxyTwo := getTestApiProxy("example-two", "foo", "/api/v2", "/")
  testApiProxyThree := getTestApiProxy("example-three", "foo", "/api", "/")
  testApiProxyFour := getTestApiProxy("example-four", "foo", "/apis", "/")
  testApiProxyFive := getTestApiProxy("example-five", "foo", "/", "/")

	assert.Nil(t, ApiProxyStore().Set(testApiProxyOne))
  assert.Nil(t, ApiProxyStore().Set(testApiProxyTwo))
  assert.Nil(t, ApiProxyStore().Set(testApiProxyThree))
  assert.Nil(t, ApiProxyStore().Set(testApiProxyFour))
  assert.Nil(t, ApiProxyStore().Set(testApiProxyFive))

	result, _ := ApiProxyStore().Get("/")
	assert.Equal(t, testApiProxyFive, result)
  result, _ = ApiProxyStore().Get("")
	assert.Equal(t, testApiProxyFive, result)
  result, _ = ApiProxyStore().Get("foo")
	assert.Equal(t, testApiProxyFive, result)
	result, _ = ApiProxyStore().Get("api/v1")
	assert.Equal(t, testApiProxyOne, result)
  result, _ = ApiProxyStore().Get("api/v2/foo")
	assert.Equal(t, testApiProxyTwo, result)
  result, _ = ApiProxyStore().Get("apis/v2/foo")
	assert.Equal(t, testApiProxyFour, result)
}

func TestApiProxyStoreDelete(t *testing.T) {
  defer ApiProxyStore().Clear()

  testApiProxyOne := getTestApiProxy("example-one", "foo", "/api/v1", "/")
  testApiProxyTwo := getTestApiProxy("example-two", "foo", "/api/v2", "/")
  testApiProxyThree := getTestApiProxy("example-three", "foo", "/api", "/")
  testApiProxyFour := getTestApiProxy("example-four", "foo", "/apis", "/")
  testApiProxyFive := getTestApiProxy("example-five", "foo", "/", "/")

	assert.Nil(t, ApiProxyStore().Set(testApiProxyOne))
  assert.Nil(t, ApiProxyStore().Set(testApiProxyTwo))
  assert.Nil(t, ApiProxyStore().Set(testApiProxyThree))
  assert.Nil(t, ApiProxyStore().Set(testApiProxyFour))

	result, _ := ApiProxyStore().Delete(nil)
	assert.Nil(t, result)
	result, _ = ApiProxyStore().Delete(testApiProxyFive)
	assert.Nil(t, result)

  assert.Nil(t, ApiProxyStore().Set(testApiProxyFive))

	result, _ = ApiProxyStore().Delete(testApiProxyOne)
	assert.Equal(t, testApiProxyOne, result)
  assert.Nil(t, apiProxyStore.proxyTree.children["api"].children["v1"])
  assert.Equal(t, 1, len(apiProxyStore.proxyTree.children["api"].children))

  result, _ = ApiProxyStore().Delete(testApiProxyTwo)
	assert.Equal(t, testApiProxyTwo, result)
  assert.Nil(t, apiProxyStore.proxyTree.children["api"].children["v1"])
  assert.Zero(t, len(apiProxyStore.proxyTree.children["api"].children))

  result, _ = ApiProxyStore().Delete(testApiProxyThree)
	assert.Equal(t, testApiProxyThree, result)
  assert.Nil(t, apiProxyStore.proxyTree.children["api"])
  assert.Equal(t, 1, len(apiProxyStore.proxyTree.children))

  result, _ = ApiProxyStore().Delete(testApiProxyFour)
	assert.Equal(t, testApiProxyFour, result)
  assert.Nil(t, apiProxyStore.proxyTree.children["apis"])
  assert.Zero(t, len(apiProxyStore.proxyTree.children))

  result, _ = ApiProxyStore().Delete(testApiProxyFive)
	assert.Equal(t, testApiProxyFive, result)
  assert.Zero(t, len(apiProxyStore.proxyTree.children))
}

func TestNormalizeProxyPaths(t *testing.T) {
  testApiProxyOne := getTestApiProxy("example-one", "foo", "///foo/bar///car", "foo///bar//car")
  testApiProxyTwo := getTestApiProxy("example-one", "foo", "", "///")

	normalizeProxyPaths(testApiProxyOne)
	normalizeProxyPaths(testApiProxyTwo)

	assert.Equal(t, testApiProxyOne.Spec.Source.Path, "/foo/bar/car")
	assert.Equal(t, testApiProxyOne.Spec.Target.Path, "/foo/bar/car")
	assert.Equal(t, testApiProxyTwo.Spec.Source.Path, "/")
	assert.Equal(t, testApiProxyTwo.Spec.Target.Path, "/")
}

func getTestApiProxy(name, namespace, sourcePath, targetPath string) *v2.ApiProxy {
  return &v2.ApiProxy{
    ObjectMeta: metav1.ObjectMeta{
      Name:      name,
      Namespace: namespace,
    },
    Spec: v2.ApiProxySpec{
      Source: v2.Source{
        Path:   sourcePath,
      },
      Target: v2.Target{
        Path: targetPath,
        Backend: v2.Backend{
          Service: v2.Service{
            Name:      "my-service",
            Port:      8080,
          },
        },
        SSL: v2.SSL{
          SecretName: "mySecret",
        },
      },
      Plugins: []v2.Plugin{
        {
          Name:    "apikey",
          Version: "v1.0.0",
        },
      },
    },
  }
}
