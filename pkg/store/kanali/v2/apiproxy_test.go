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

	"github.com/northwesternmutual/kanali/test/builder"
)

func TestApiProxyStore(t *testing.T) {
	_, ok := ApiProxyStore().(ApiProxyStoreInterface)
	assert.True(t, ok)
}

func TestApiProxyStoreSet(t *testing.T) {
	defer ApiProxyStore().Clear()

	testApiProxyOne := builder.NewApiProxy("example-one", "foo").WithSourcePath("/api/v1").NewOrDie()
	testApiProxyTwo := builder.NewApiProxy("example-two", "foo").WithSourcePath("/api/v2").NewOrDie()
	testApiProxyThree := builder.NewApiProxy("example-three", "foo").WithSourcePath("/api").NewOrDie()
	testApiProxyFour := builder.NewApiProxy("example-four", "foo").WithSourcePath("/apis").NewOrDie()
	testApiProxyFive := builder.NewApiProxy("example-five", "foo").WithSourcePath("/").NewOrDie()
	testApiProxySix := builder.NewApiProxy("example-six", "foo").WithSourcePath("/api/v3").WithSourceHost("foo.bar.com").NewOrDie()
	testApiProxySeven := builder.NewApiProxy("example-seven", "foo").WithSourcePath("/api/v3").WithSourceHost("bar.foo.com").NewOrDie()

	ApiProxyStore().Set(testApiProxyOne)
	ApiProxyStore().Set(testApiProxyTwo)
	ApiProxyStore().Set(testApiProxyThree)
	ApiProxyStore().Set(testApiProxyFour)
	ApiProxyStore().Set(testApiProxyFive)
	ApiProxyStore().Set(testApiProxySix)
	ApiProxyStore().Set(testApiProxySeven)

	assert.Equal(t, testApiProxyOne, apiProxyStore.proxyTree.children["api"].children["v1"].value.global)
	assert.Equal(t, testApiProxyTwo, apiProxyStore.proxyTree.children["api"].children["v2"].value.global)
	assert.Equal(t, testApiProxyThree, apiProxyStore.proxyTree.children["api"].value.global)
	assert.Equal(t, testApiProxyFour, apiProxyStore.proxyTree.children["apis"].value.global)
	assert.Equal(t, testApiProxyFive, apiProxyStore.proxyTree.value.global)
	assert.Nil(t, apiProxyStore.proxyTree.children["api"].children["v3"].value.global)
	assert.Equal(t, testApiProxySix, apiProxyStore.proxyTree.children["api"].children["v3"].value.vhosts["foo.bar.com"])
	assert.Equal(t, testApiProxySeven, apiProxyStore.proxyTree.children["api"].children["v3"].value.vhosts["bar.foo.com"])
}

func TestApiProxyStoreClear(t *testing.T) {
	ApiProxyStore().Set(builder.NewApiProxy("example-one", "foo").WithSourcePath("/api/v1").NewOrDie())
	assert.Equal(t, 1, len(apiProxyStore.proxyTree.children))
	ApiProxyStore().Clear()
	assert.Equal(t, 0, len(apiProxyStore.proxyTree.children))
}

func TestApiProxyStoreUpdate(t *testing.T) {
	defer ApiProxyStore().Clear()

	testApiProxyOneOld := builder.NewApiProxy("example-one", "foo").WithSourcePath("/api/v1").WithTargetPath("/foo").NewOrDie()
	testApiProxyTwoOld := builder.NewApiProxy("example-two", "foo").WithSourcePath("/apis").WithTargetPath("/").NewOrDie()
	testApiProxyThree := builder.NewApiProxy("example-three", "foo").WithSourcePath("/api/v1").WithTargetPath("/").NewOrDie()

	ApiProxyStore().Set(testApiProxyOneOld)
	ApiProxyStore().Set(testApiProxyTwoOld)

	testApiProxyOneNew := builder.NewApiProxy("example-one", "foo").WithSourcePath("/api/v1").WithTargetPath("/").NewOrDie()
	testApiProxyTwoNew := builder.NewApiProxy("example-two", "foo").WithSourcePath("/").WithTargetPath("/").NewOrDie()

	ApiProxyStore().Update(testApiProxyOneOld, testApiProxyOneNew)
	ApiProxyStore().Update(testApiProxyTwoOld, testApiProxyTwoNew)

	assert.Equal(t, testApiProxyOneNew, apiProxyStore.proxyTree.children["api"].children["v1"].value.global)
	assert.Equal(t, testApiProxyTwoNew, apiProxyStore.proxyTree.children[""].value.global)
	assert.Nil(t, apiProxyStore.proxyTree.children["apis"])

	ApiProxyStore().Clear()
	assert.True(t, ApiProxyStore().IsEmpty())
	ApiProxyStore().Set(testApiProxyOneOld)
	ApiProxyStore().Update(testApiProxyOneOld, testApiProxyThree)
	assert.Equal(t, testApiProxyOneOld, apiProxyStore.proxyTree.children["api"].children["v1"].value.global)
}

func TestApiProxyStoreIsEmpty(t *testing.T) {
	defer ApiProxyStore().Clear()

	assert.True(t, ApiProxyStore().IsEmpty())
	ApiProxyStore().Set(builder.NewApiProxy("example-one", "foo").WithSourcePath("/api/v1").NewOrDie())
	assert.False(t, ApiProxyStore().IsEmpty())
}

func TestApiProxyStoreGet(t *testing.T) {
	defer ApiProxyStore().Clear()

	testApiProxyOne := builder.NewApiProxy("example-one", "foo").WithSourcePath("/api/v1").NewOrDie()
	testApiProxyTwo := builder.NewApiProxy("example-two", "foo").WithSourcePath("/api/v2").NewOrDie()
	testApiProxyThree := builder.NewApiProxy("example-three", "foo").WithSourcePath("/api").NewOrDie()
	testApiProxyFour := builder.NewApiProxy("example-four", "foo").WithSourcePath("/apis").NewOrDie()
	testApiProxyFive := builder.NewApiProxy("example-five", "foo").WithSourcePath("/").NewOrDie()
	testApiProxySix := builder.NewApiProxy("example-six", "foo").WithSourcePath("/api/v3").WithSourceHost("foo.bar.com").NewOrDie()

	assert.Nil(t, ApiProxyStore().Get("/", ""))

	ApiProxyStore().Set(testApiProxyOne)
	ApiProxyStore().Set(testApiProxyTwo)
	ApiProxyStore().Set(testApiProxyThree)
	ApiProxyStore().Set(testApiProxyFour)
	ApiProxyStore().Set(testApiProxyFive)
	ApiProxyStore().Set(testApiProxySix)

	assert.Equal(t, testApiProxyFive, ApiProxyStore().Get("/", ""))
	assert.Equal(t, testApiProxyFive, ApiProxyStore().Get("", ""))
	assert.Equal(t, testApiProxyFive, ApiProxyStore().Get("foo", ""))
	assert.Equal(t, testApiProxyOne, ApiProxyStore().Get("api/v1", ""))
	assert.Equal(t, testApiProxyTwo, ApiProxyStore().Get("api/v2/foo", ""))
	assert.Equal(t, testApiProxyFour, ApiProxyStore().Get("apis/v2/foo", ""))

	ApiProxyStore().Clear()
	testApiProxySeven := builder.NewApiProxy("example-seven", "foo").WithSourcePath("/").WithSourceHost("foo.bar.com").NewOrDie()
	ApiProxyStore().Set(testApiProxySeven)
	assert.Equal(t, testApiProxySeven, ApiProxyStore().Get("/frank", "foo.bar.com"))
}

func TestApiProxyStoreDelete(t *testing.T) {
	defer ApiProxyStore().Clear()

	testApiProxyOne := builder.NewApiProxy("example-one", "foo").WithSourcePath("/api/v1").NewOrDie()
	testApiProxyTwo := builder.NewApiProxy("example-two", "foo").WithSourcePath("/api/v2").NewOrDie()
	testApiProxyThree := builder.NewApiProxy("example-three", "foo").WithSourcePath("/api").NewOrDie()
	testApiProxyFour := builder.NewApiProxy("example-four", "foo").WithSourcePath("/apis").NewOrDie()
	testApiProxyFive := builder.NewApiProxy("example-five", "foo").WithSourcePath("/").NewOrDie()
	testApiProxySix := builder.NewApiProxy("example-six", "foo").WithSourcePath("/api/v3").WithSourceHost("foo.bar.com").NewOrDie()
	testApiProxySeven := builder.NewApiProxy("example-seven", "foo").WithSourcePath("/api/v3").WithSourceHost("bar.foo.com").NewOrDie()

	ApiProxyStore().Set(testApiProxyOne)
	ApiProxyStore().Set(testApiProxyTwo)
	ApiProxyStore().Set(testApiProxyThree)
	ApiProxyStore().Set(testApiProxyFour)

	assert.Nil(t, ApiProxyStore().Delete(nil))
	assert.Nil(t, ApiProxyStore().Delete(testApiProxyFive))

	ApiProxyStore().Set(testApiProxyFive)

	assert.Equal(t, testApiProxyOne, ApiProxyStore().Delete(testApiProxyOne))
	assert.Nil(t, apiProxyStore.proxyTree.children["api"].children["v1"])
	assert.Equal(t, 1, len(apiProxyStore.proxyTree.children["api"].children))

	assert.Equal(t, testApiProxyTwo, ApiProxyStore().Delete(testApiProxyTwo))
	assert.Nil(t, apiProxyStore.proxyTree.children["api"].children["v1"])
	assert.Zero(t, len(apiProxyStore.proxyTree.children["api"].children))

	assert.Equal(t, testApiProxyThree, ApiProxyStore().Delete(testApiProxyThree))
	assert.Nil(t, apiProxyStore.proxyTree.children["api"])
	assert.Equal(t, 1, len(apiProxyStore.proxyTree.children))

	assert.Equal(t, testApiProxyFour, ApiProxyStore().Delete(testApiProxyFour))
	assert.Nil(t, apiProxyStore.proxyTree.children["apis"])
	assert.Zero(t, len(apiProxyStore.proxyTree.children))

	assert.Equal(t, testApiProxyFive, ApiProxyStore().Delete(testApiProxyFive))
	assert.Zero(t, len(apiProxyStore.proxyTree.children))

	ApiProxyStore().Set(testApiProxySix)
	ApiProxyStore().Set(testApiProxySeven)
	assert.Equal(t, testApiProxySix, ApiProxyStore().Delete(testApiProxySix))
	assert.Equal(t, 1, len(apiProxyStore.proxyTree.children["api"].children["v3"].value.vhosts))
	assert.Equal(t, testApiProxySeven, ApiProxyStore().Delete(testApiProxySeven))
	assert.Nil(t, apiProxyStore.proxyTree.children["api"])
}

func TestNormalizeProxyPaths(t *testing.T) {
	testApiProxyOne := builder.NewApiProxy("example-one", "foo").WithSourcePath("///foo/bar///car").WithTargetPath("foo///bar//car").NewOrDie()
	testApiProxyTwo := builder.NewApiProxy("example-one", "foo").WithSourcePath("").WithTargetPath("///").NewOrDie()

	normalizeProxyPaths(testApiProxyOne)
	normalizeProxyPaths(testApiProxyTwo)

	assert.Equal(t, testApiProxyOne.Spec.Source.Path, "/foo/bar/car")
	assert.Equal(t, testApiProxyOne.Spec.Target.Path, "/foo/bar/car")
	assert.Equal(t, testApiProxyTwo.Spec.Source.Path, "/")
	assert.Equal(t, testApiProxyTwo.Spec.Target.Path, "/")
}
