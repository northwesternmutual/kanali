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

	"github.com/northwesternmutual/kanali/pkg/apis/kanali.io/v2"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestApiKeyStore(t *testing.T) {
	_, ok := ApiKeyStore().(ApiKeyStoreInterface)
	assert.True(t, ok)
}

func TestApiKeyStoreSet(t *testing.T) {
	defer ApiKeyStore().Clear()

	apiKeyOne := getTestApiKey("example-one", "foo", v2.Revision{Data: "abc123", Status: "active"})
	apiKeyTwo := getTestApiKey("example-two", "foo", v2.Revision{Data: "def456", Status: "active"}, v2.Revision{Data: "ghi789", Status: "active"})

	ApiKeyStore().Set(apiKeyOne)
	ApiKeyStore().Set(apiKeyTwo)

	assert.Equal(t, *apiKeyOne, apiKeyStore.keyMap["abc123"])
	assert.Equal(t, *apiKeyTwo, apiKeyStore.keyMap["def456"])
	assert.Equal(t, *apiKeyTwo, apiKeyStore.keyMap["ghi789"])
}

func TestApiKeyStoreUpdate(t *testing.T) {
	defer ApiKeyStore().Clear()

	apiKeyOneOld := getTestApiKey("example-one", "foo", v2.Revision{Data: "mno345", Status: "active"}, v2.Revision{Data: "abc123", Status: "active"})
	apiKeyOneNew := getTestApiKey("example-one", "foo", v2.Revision{Data: "abc123", Status: "inactive"})

	apiKeyTwoOld := getTestApiKey("example-two", "foo", v2.Revision{Data: "def456", Status: "active"}, v2.Revision{Data: "ghi789", Status: "active"})
	apiKeyTwoNew := getTestApiKey("example-two", "foo", v2.Revision{Data: "jkl012", Status: "active"})

	apiKeyThreeOld := getTestApiKey("example-three", "foo", v2.Revision{Data: "cba321", Status: "active"})
	apiKeyThreeNew := getTestApiKey("example-three", "foo", v2.Revision{Data: "abc123", Status: "active"})

	ApiKeyStore().Update(apiKeyOneOld, apiKeyOneNew)
	ApiKeyStore().Update(apiKeyTwoOld, apiKeyTwoNew)

	assert.Equal(t, *apiKeyOneNew, apiKeyStore.keyMap["abc123"])
	assert.Equal(t, *apiKeyTwoNew, apiKeyStore.keyMap["jkl012"])
	assert.Equal(t, v2.ApiKey{}, apiKeyStore.keyMap["def456"])
	assert.Equal(t, v2.ApiKey{}, apiKeyStore.keyMap["ghi789"])
	assert.Equal(t, v2.ApiKey{}, apiKeyStore.keyMap["mno345"])

	ApiKeyStore().Clear()
	ApiKeyStore().Update(apiKeyThreeOld, apiKeyThreeNew)
	assert.Equal(t, v2.ApiKey{}, apiKeyStore.keyMap["cba321"])
	assert.Equal(t, *apiKeyThreeNew, apiKeyStore.keyMap["abc123"])
}

func TestApiKeyStoreClear(t *testing.T) {
	defer ApiKeyStore().Clear()

	apiKeyOne := getTestApiKey("example-one", "foo", v2.Revision{Data: "abc123", Status: "active"})

	ApiKeyStore().Set(apiKeyOne)
	assert.Equal(t, 1, len(apiKeyStore.keyMap))
	ApiKeyStore().Clear()
	assert.Equal(t, 0, len(apiKeyStore.keyMap))
}

func TestApiKeyStoreIsEmpty(t *testing.T) {
	defer ApiKeyStore().Clear()

	apiKeyOne := getTestApiKey("example-one", "foo", v2.Revision{Data: "abc123", Status: "active"})

	assert.True(t, ApiKeyStore().IsEmpty())
	ApiKeyStore().Set(apiKeyOne)
	assert.False(t, ApiKeyStore().IsEmpty())
	ApiKeyStore().Clear()
	assert.True(t, ApiKeyStore().IsEmpty())
}

func TestApiKeyStoreGet(t *testing.T) {
	defer ApiKeyStore().Clear()

	apiKeyOne := getTestApiKey("example-one", "foo", v2.Revision{Data: "abc123", Status: "active"})
	apiKeyTwo := getTestApiKey("example-two", "foo", v2.Revision{Data: "def456", Status: "active"}, v2.Revision{Data: "ghi789", Status: "active"})

	ApiKeyStore().Set(apiKeyOne)
	ApiKeyStore().Set(apiKeyTwo)

	assert.Nil(t, ApiKeyStore().Get(""))
	assert.Equal(t, apiKeyOne, ApiKeyStore().Get("abc123"))
	assert.Equal(t, apiKeyTwo, ApiKeyStore().Get("def456"))
	assert.Equal(t, apiKeyTwo, ApiKeyStore().Get("ghi789"))
}

func TestAPIKeyDelete(t *testing.T) {
	defer ApiKeyStore().Clear()

	apiKeyOne := getTestApiKey("example-one", "foo", v2.Revision{Data: "abc123", Status: "active"})
	apiKeyTwo := getTestApiKey("example-two", "foo", v2.Revision{Data: "def456", Status: "active"}, v2.Revision{Data: "ghi789", Status: "active"})

	ApiKeyStore().Set(apiKeyOne)
	ApiKeyStore().Set(apiKeyTwo)

	assert.Nil(t, ApiKeyStore().Delete(nil))
	assert.Equal(t, apiKeyOne, ApiKeyStore().Delete(apiKeyOne))
	assert.Equal(t, apiKeyTwo, ApiKeyStore().Delete(apiKeyTwo))
	assert.True(t, ApiKeyStore().IsEmpty())
	assert.Nil(t, ApiKeyStore().Delete(getTestApiKey("example-one", "foo")))
	assert.Nil(t, ApiKeyStore().Delete(getTestApiKey("example-one", "foo", v2.Revision{Data: "abc123", Status: "active"})))
}

func getTestApiKey(name, namespace string, revisions ...v2.Revision) *v2.ApiKey {
	return &v2.ApiKey{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: v2.ApiKeySpec{
			Revisions: revisions,
		},
	}
}
