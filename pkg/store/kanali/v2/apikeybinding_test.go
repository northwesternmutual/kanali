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

func TestApiKeyBindingStore(t *testing.T) {
	_, ok := ApiKeyBindingStore().(ApiKeyBindingStoreInterface)
	assert.True(t, ok)
}

func TestApiKeyBindingStoreSet(t *testing.T) {
	defer ApiKeyBindingStore().Clear()

	testApiKeyBindingOne := builder.NewApiKeyBinding("example-one", "foo").WithKeys(
		builder.NewKeyAccess("example-one").WithDefaultRule(
			builder.NewRule().WithGranular("POST", "PUT").NewOrDie(),
		).WithSubpaths(
			builder.NewPathBuilder("/foo", builder.NewRule().WithGlobal().NewOrDie()).NewOrDie(),
			builder.NewPathBuilder("/foo/bar", builder.NewRule().WithGlobal().NewOrDie()).NewOrDie(),
		).NewOrDie(),
	).NewOrDie()

	ApiKeyBindingStore().Set(testApiKeyBindingOne)

	assert.Equal(t, 1, len(apiKeyBindingStore.apiKeyBindingMap["foo"]))
	assert.Equal(t, 1, len(apiKeyBindingStore.apiKeyBindingMap["foo"]["example-one"]))
	assert.Equal(t, testApiKeyBindingOne.Spec.Keys[0], apiKeyBindingStore.apiKeyBindingMap["foo"]["example-one"]["example-one"])
}

func TestApiKeyBindingStoreClear(t *testing.T) {
	defer ApiKeyBindingStore().Clear()

	assert.Equal(t, 0, len(apiKeyBindingStore.apiKeyBindingMap))

	testApiKeyBindingOne := builder.NewApiKeyBinding("example-one", "foo").WithKeys(
		builder.NewKeyAccess("example-one").WithDefaultRule(
			builder.NewRule().WithGranular("POST", "PUT").NewOrDie(),
		).WithSubpaths(
			builder.NewPathBuilder("/foo", builder.NewRule().WithGlobal().NewOrDie()).NewOrDie(),
		).NewOrDie(),
	).NewOrDie()

	ApiKeyBindingStore().Set(testApiKeyBindingOne)

	assert.Equal(t, 1, len(apiKeyBindingStore.apiKeyBindingMap))
	ApiKeyBindingStore().Clear()
	assert.Equal(t, 0, len(apiKeyBindingStore.apiKeyBindingMap))
}

func TestAPIKeyBindingUpdate(t *testing.T) {
	defer ApiKeyBindingStore().Clear()

	testApiKeyBindingOld := builder.NewApiKeyBinding("example-one", "foo").WithKeys(
		builder.NewKeyAccess("example-one").WithDefaultRule(
			builder.NewRule().WithGranular("POST", "PUT").NewOrDie(),
		).WithSubpaths(
			builder.NewPathBuilder("/foo", builder.NewRule().WithGranular("GET", "DELETE").NewOrDie()).NewOrDie(),
			builder.NewPathBuilder("/foo/bar", builder.NewRule().WithGlobal().NewOrDie()).NewOrDie(),
		).NewOrDie(),
	).NewOrDie()

	testApiKeyBindingNew := builder.NewApiKeyBinding("example-one", "foo").WithKeys(
		builder.NewKeyAccess("example-one").WithDefaultRule(
			builder.NewRule().WithGranular("POST", "PUT").NewOrDie(),
		).WithSubpaths(
			builder.NewPathBuilder("/foo/bar", builder.NewRule().WithGlobal().NewOrDie()).NewOrDie(),
		).NewOrDie(),
	).NewOrDie()

	ApiKeyBindingStore().Set(testApiKeyBindingOld)
	ApiKeyBindingStore().Update(testApiKeyBindingOld, testApiKeyBindingNew)
}

func TestApiKeyBindingIsEmpty(t *testing.T) {
	defer ApiKeyBindingStore().Clear()
	assert.True(t, ApiKeyBindingStore().IsEmpty())

	ApiKeyBindingStore().Set(
		builder.NewApiKeyBinding("foo", "bar").WithKeys(
			builder.NewKeyAccess("frank").WithDefaultRule(
				builder.NewRule().WithGlobal().NewOrDie(),
			).NewOrDie(),
		).NewOrDie(),
	)

	assert.False(t, ApiKeyBindingStore().IsEmpty())
	ApiKeyBindingStore().Clear()
	assert.True(t, ApiKeyBindingStore().IsEmpty())
}

func TestContains(t *testing.T) {
	defer ApiKeyBindingStore().Clear()
	assert.False(t, ApiKeyBindingStore().Contains("bar", "foo"))

	ApiKeyBindingStore().Set(
		builder.NewApiKeyBinding("foo", "bar").WithKeys(
			builder.NewKeyAccess("frank").WithDefaultRule(
				builder.NewRule().WithGlobal().NewOrDie(),
			).NewOrDie(),
		).NewOrDie(),
	)

	assert.True(t, ApiKeyBindingStore().Contains("bar", "foo"))
	ApiKeyBindingStore().Clear()
	assert.False(t, ApiKeyBindingStore().Contains("bar", "foo"))
}

func TestGetHightestPriorityRule(t *testing.T) {
	defer ApiKeyBindingStore().Clear()
	assert.Nil(t, ApiKeyBindingStore().GetHightestPriorityRule("foo", "bar", "car", ""))

	globalRule := builder.NewRule().WithGlobal().NewOrDie()
	nothingRule := builder.NewRule().NewOrDie()

	ApiKeyBindingStore().Set(
		builder.NewApiKeyBinding("foo", "bar").WithKeys(
			builder.NewKeyAccess("frank").WithDefaultRule(globalRule).NewOrDie(),
		).NewOrDie(),
	)
	assert.Equal(t, globalRule, *ApiKeyBindingStore().GetHightestPriorityRule("bar", "foo", "frank", ""))

	ApiKeyBindingStore().Set(
		builder.NewApiKeyBinding("foo", "bar").WithKeys(
			builder.NewKeyAccess("frank").WithDefaultRule(nothingRule).NewOrDie(),
		).NewOrDie(),
	)
	assert.Equal(t, nothingRule, *ApiKeyBindingStore().GetHightestPriorityRule("bar", "foo", "frank", ""))

	ApiKeyBindingStore().Set(
		builder.NewApiKeyBinding("foo", "bar").WithKeys(
			builder.NewKeyAccess("frank").WithDefaultRule(nothingRule).WithSubpaths(
				builder.NewPathBuilder("/foo", globalRule).NewOrDie(),
			).NewOrDie(),
		).NewOrDie(),
	)
	assert.Equal(t, globalRule, *ApiKeyBindingStore().GetHightestPriorityRule("bar", "foo", "frank", "/foo"))
	assert.Equal(t, globalRule, *ApiKeyBindingStore().GetHightestPriorityRule("bar", "foo", "frank", "/foo/bar"))
	assert.Equal(t, nothingRule, *ApiKeyBindingStore().GetHightestPriorityRule("bar", "foo", "frank", "/bar"))

	ApiKeyBindingStore().Set(
		builder.NewApiKeyBinding("foo", "bar").WithKeys(
			builder.NewKeyAccess("frank").WithDefaultRule(nothingRule).WithSubpaths(
				builder.NewPathBuilder("/foo$", globalRule).NewOrDie(),
				builder.NewPathBuilder("/foo*", nothingRule).NewOrDie(),
			).NewOrDie(),
		).NewOrDie(),
	)
	assert.Equal(t, globalRule, *ApiKeyBindingStore().GetHightestPriorityRule("bar", "foo", "frank", "/foo"))
	assert.Equal(t, nothingRule, *ApiKeyBindingStore().GetHightestPriorityRule("bar", "foo", "frank", "/fooo"))
}

func TestApiKeyBindingDelete(t *testing.T) {
	defer ApiKeyBindingStore().Clear()

	globalRule := builder.NewRule().WithGlobal().NewOrDie()

	apikeybindingOne := builder.NewApiKeyBinding("foo", "bar").WithKeys(
		builder.NewKeyAccess("frank").WithDefaultRule(globalRule).NewOrDie(),
	).NewOrDie()
	apikeybindingTwo := builder.NewApiKeyBinding("car", "bar").WithKeys(
		builder.NewKeyAccess("frank").WithDefaultRule(globalRule).NewOrDie(),
	).NewOrDie()

	assert.Nil(t, ApiKeyBindingStore().Delete(nil))
	assert.NotNil(t, ApiKeyBindingStore().Delete(apikeybindingOne))
	ApiKeyBindingStore().Set(apikeybindingOne)
	ApiKeyBindingStore().Set(apikeybindingTwo)
	assert.Equal(t, 2, len(apiKeyBindingStore.apiKeyBindingMap["bar"]))
	assert.Equal(t, 1, len(apiKeyBindingStore.apiKeyBindingMap))
	assert.Nil(t, ApiKeyBindingStore().Delete(apikeybindingTwo))
	assert.Equal(t, 1, len(apiKeyBindingStore.apiKeyBindingMap))
	assert.Equal(t, 1, len(apiKeyBindingStore.apiKeyBindingMap["bar"]))
	assert.Nil(t, ApiKeyBindingStore().Delete(apikeybindingOne))
	assert.Equal(t, 0, len(apiKeyBindingStore.apiKeyBindingMap))
}
