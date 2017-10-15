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
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestGetMockResponseStore(t *testing.T) {
	MockResponseStore.Clear()
	assert.Equal(t, 0, len(MockResponseStore.mockRespTree))

	v := MockResponseFactory{}
	var i interface{} = &v
	_, ok := i.(Store)
	assert.True(t, ok)
}

func TestIsValidHTTPMethod(t *testing.T) {
	assert.False(t, isValidHTTPMethod("foo"))
	assert.True(t, isValidHTTPMethod("GET"))
	assert.True(t, isValidHTTPMethod("get"))
	assert.True(t, isValidHTTPMethod("POST"))
	assert.True(t, isValidHTTPMethod("post"))
}

func TestMockResponseSet(t *testing.T) {
	cm := getTestConfigMaps()

	MockResponseStore.Clear()

	assert.Equal(t, MockResponseStore.Set("hi").Error(), "obj was not a ConfigMap")
	assert.Nil(t, MockResponseStore.Set(v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "cm-one",
			Namespace: "foo",
		},
		Data: map[string]string{
			"foo": "bar",
		},
	}))
	assert.Nil(t, MockResponseStore.Set(v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "cm-one",
			Namespace: "foo",
		},
		Data: map[string]string{
			"response": "bad",
		},
	}))
	MockResponseStore.Set(cm[0])
	assert.Equal(t, len(MockResponseStore.mockRespTree), 1)
	assert.Equal(t, len(MockResponseStore.mockRespTree["foo"]), 1)
	assert.Equal(t, len(MockResponseStore.mockRespTree["foo"]["cm-one"]), 1)

	MockResponseStore.Set(cm[1])
	assert.Equal(t, len(MockResponseStore.mockRespTree), 1)
	assert.Equal(t, len(MockResponseStore.mockRespTree["foo"]), 2)
	assert.Equal(t, len(MockResponseStore.mockRespTree["foo"]["cm-two"]["GET"].Children), 1)
	assert.Equal(t, MockResponseStore.mockRespTree["foo"]["cm-two"]["GET"].Children["bar"].Value, &Route{
		Route:  "/bar",
		Code:   200,
		Method: "GET",
		Body:   "{\"foo\": \"bar\"}",
	})
	assert.Equal(t, MockResponseStore.mockRespTree["foo"]["cm-two"]["GET"].Value, &Route{
		Route:  "/",
		Code:   200,
		Method: "GET",
		Body:   "{\"frank\": \"greco\"}",
	})

	MockResponseStore.Set(cm[2])
	assert.Equal(t, len(MockResponseStore.mockRespTree), 1)
	assert.Equal(t, len(MockResponseStore.mockRespTree["foo"]), 3)
	assert.Equal(t, len(MockResponseStore.mockRespTree["foo"]["cm-three"]["GET"].Children), 1)
	assert.Equal(t, MockResponseStore.mockRespTree["foo"]["cm-three"]["GET"].Value, &Route{
		Route:  "",
		Code:   200,
		Method: "GET",
		Body:   "{\"foo\": \"bar\"}",
	})
	assert.Equal(t, MockResponseStore.mockRespTree["foo"]["cm-three"]["GET"].Children["foo"].Value, &Route{
		Route:  "foo",
		Code:   200,
		Method: "GET",
		Body:   "{\"foo\": \"bar\"}",
	})
	assert.Equal(t, MockResponseStore.mockRespTree["foo"]["cm-three"]["POST"].Children["foo"].Children["bar"].Children["car"].Value, &Route{
		Route:  "/foo/bar/car",
		Code:   200,
		Method: "POST",
		Body:   "{\"foo\": \"bar\"}",
	})
	assert.Equal(t, len(MockResponseStore.mockRespTree["foo"]["cm-three"]["GET"].Children["foo"].Children), 0)
	assert.Equal(t, len(MockResponseStore.mockRespTree["foo"]["cm-three"]["POST"].Children["foo"].Children), 1)

	MockResponseStore.Set(cm[3])
	assert.Equal(t, len(MockResponseStore.mockRespTree), 2)
	assert.Equal(t, len(MockResponseStore.mockRespTree["bar"]), 1)
	assert.Equal(t, len(MockResponseStore.mockRespTree["bar"]["cm-four"]), 1)

}

func TestMockResponseUpdate(t *testing.T) {
	cm := getTestConfigMaps()

	MockResponseStore.Clear()

	assert.Equal(t, MockResponseStore.Update("hi", "hi").Error(), "obj was not a ConfigMap")
	configMapOne := v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "cm-one",
			Namespace: "foo",
		},
		Data: map[string]string{
			"foo": "bar",
		},
	}
	configMapTwo := v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "cm-one",
			Namespace: "foo",
		},
		Data: map[string]string{
			"response": "bad",
		},
	}
	assert.Nil(t, MockResponseStore.Update(configMapOne, configMapOne))
	assert.Nil(t, MockResponseStore.Update(configMapTwo, configMapTwo))
	MockResponseStore.Update(cm[0], cm[0])
	assert.Equal(t, len(MockResponseStore.mockRespTree), 1)
	assert.Equal(t, len(MockResponseStore.mockRespTree["foo"]), 1)
	assert.Equal(t, len(MockResponseStore.mockRespTree["foo"]["cm-one"]), 1)

	MockResponseStore.Update(cm[1], cm[1])
	assert.Equal(t, len(MockResponseStore.mockRespTree), 1)
	assert.Equal(t, len(MockResponseStore.mockRespTree["foo"]), 2)
	assert.Equal(t, len(MockResponseStore.mockRespTree["foo"]["cm-two"]["GET"].Children), 1)
	assert.Equal(t, MockResponseStore.mockRespTree["foo"]["cm-two"]["GET"].Children["bar"].Value, &Route{
		Route:  "/bar",
		Code:   200,
		Method: "GET",
		Body:   "{\"foo\": \"bar\"}",
	})
	assert.Equal(t, MockResponseStore.mockRespTree["foo"]["cm-two"]["GET"].Value, &Route{
		Route:  "/",
		Code:   200,
		Method: "GET",
		Body:   "{\"frank\": \"greco\"}",
	})

	MockResponseStore.Update(cm[2], cm[2])
	assert.Equal(t, len(MockResponseStore.mockRespTree), 1)
	assert.Equal(t, len(MockResponseStore.mockRespTree["foo"]), 3)
	assert.Equal(t, len(MockResponseStore.mockRespTree["foo"]["cm-three"]["GET"].Children), 1)
	assert.Equal(t, MockResponseStore.mockRespTree["foo"]["cm-three"]["GET"].Value, &Route{
		Route:  "",
		Code:   200,
		Method: "GET",
		Body:   "{\"foo\": \"bar\"}",
	})
	assert.Equal(t, MockResponseStore.mockRespTree["foo"]["cm-three"]["GET"].Children["foo"].Value, &Route{
		Route:  "foo",
		Code:   200,
		Method: "GET",
		Body:   "{\"foo\": \"bar\"}",
	})
	assert.Equal(t, MockResponseStore.mockRespTree["foo"]["cm-three"]["POST"].Children["foo"].Children["bar"].Children["car"].Value, &Route{
		Route:  "/foo/bar/car",
		Code:   200,
		Method: "POST",
		Body:   "{\"foo\": \"bar\"}",
	})
	assert.Equal(t, len(MockResponseStore.mockRespTree["foo"]["cm-three"]["GET"].Children["foo"].Children), 0)
	assert.Equal(t, len(MockResponseStore.mockRespTree["foo"]["cm-three"]["POST"].Children["foo"].Children), 1)

	MockResponseStore.Set(cm[3])
	assert.Equal(t, len(MockResponseStore.mockRespTree), 2)
	assert.Equal(t, len(MockResponseStore.mockRespTree["bar"]), 1)
	assert.Equal(t, len(MockResponseStore.mockRespTree["bar"]["cm-four"]), 1)

}

func TestMockResponseClear(t *testing.T) {
	cm := getTestConfigMaps()
	MockResponseStore.Clear()
	MockResponseStore.Set(cm[0])
	MockResponseStore.Set(cm[1])
	MockResponseStore.Set(cm[2])
	assert.Equal(t, len(MockResponseStore.mockRespTree), 1)
	MockResponseStore.Clear()
	assert.Equal(t, len(MockResponseStore.mockRespTree), 0)
}

func TestMockResponseGet(t *testing.T) {
	cm := getTestConfigMaps()
	MockResponseStore.Clear()
	MockResponseStore.Set(cm[0])
	MockResponseStore.Set(cm[1])
	MockResponseStore.Set(cm[2])
	MockResponseStore.Set(cm[3])
	MockResponseStore.Set(cm[4])
	_, err := MockResponseStore.Get("foo")
	assert.Equal(t, err.Error(), "four parameters expected")
	_, err = MockResponseStore.Get(5, "bar", "car", "far")
	assert.Equal(t, err.Error(), "expecting namespace")
	_, err = MockResponseStore.Get("foo", 5, "car", "far")
	assert.Equal(t, err.Error(), "expecting name")
	_, err = MockResponseStore.Get("foo", "bar", 5, "far")
	assert.Equal(t, err.Error(), "expecting path")
	_, err = MockResponseStore.Get("foo", "bar", "car", 5)
	assert.Equal(t, err.Error(), "expecting method")

	result, _ := MockResponseStore.Get("bar", "bar", "bar", "bar")
	assert.Nil(t, result)
	result, _ = MockResponseStore.Get("foo", "cm-three", "", "GET")
	assert.Equal(t, result, &Route{
		Route:  "",
		Code:   200,
		Method: "GET",
		Body:   "{\"foo\": \"bar\"}",
	})
	result, _ = MockResponseStore.Get("foo", "cm-three", "/", "GET")
	assert.Equal(t, result, &Route{
		Route:  "",
		Code:   200,
		Method: "GET",
		Body:   "{\"foo\": \"bar\"}",
	})
	result, _ = MockResponseStore.Get("foo", "cm-one", "/frank/greco", "GET")
	assert.Nil(t, result)
	result, _ = MockResponseStore.Get("foo", "cm-five", "/frank/greco", "GET")
	assert.Nil(t, result)
	result, _ = MockResponseStore.Get("foo", "cm-three", "/foo/bar/car/far", "POST")
	assert.Equal(t, result, Route{
		Route:  "/foo/bar/car",
		Code:   200,
		Method: "POST",
		Body:   "{\"foo\": \"bar\"}",
	})
	result, _ = MockResponseStore.Get("foo", "cm-three", "/foo/bar/car/far", "CUSTOM")
	assert.Nil(t, result)
}

func TestMockResponseDelete(t *testing.T) {
	cm := getTestConfigMaps()
	MockResponseStore.Clear()
	MockResponseStore.Set(cm[0])
	MockResponseStore.Set(cm[1])
	MockResponseStore.Set(cm[2])
	MockResponseStore.Set(cm[3])
	MockResponseStore.Set(cm[4])
	result, err := MockResponseStore.Delete(nil)
	assert.Nil(t, result)
	assert.Nil(t, err)
	_, err = MockResponseStore.Delete(5)
	assert.Equal(t, err.Error(), "obj was not a ConfigMap")

	result, _ = MockResponseStore.Delete(cm[0])
	assert.Nil(t, result)
	assert.Equal(t, len(MockResponseStore.mockRespTree), 2)
	assert.Nil(t, MockResponseStore.mockRespTree["foo"]["cm-one"])

	result, _ = MockResponseStore.Delete(cm[3])
	assert.Nil(t, result)
	assert.Equal(t, len(MockResponseStore.mockRespTree), 1)
	assert.Nil(t, MockResponseStore.mockRespTree["bar"])
	result, _ = MockResponseStore.Delete(v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "cm-six",
			Namespace: "foo",
		},
		Data: map[string]string{
			"response": "",
		},
	})
	assert.Nil(t, result)
}

func TestMockResponseIsEmpty(t *testing.T) {
	cm := getTestConfigMaps()
	MockResponseStore.Clear()
	MockResponseStore.Set(cm[0])
	MockResponseStore.Set(cm[1])
	MockResponseStore.Set(cm[2])
	assert.False(t, MockResponseStore.IsEmpty())
	MockResponseStore.Clear()
	assert.True(t, MockResponseStore.IsEmpty())
}

func getTestConfigMaps() []v1.ConfigMap {

	mockOne, _ := json.Marshal(mock{
		Route{
			Route:  "/foo",
			Code:   200,
			Method: "GET",
			Body:   "{\"foo\": \"bar\"}",
		},
	})

	mockTwo, _ := json.Marshal(mock{
		Route{
			Route:  "/bar",
			Code:   200,
			Method: "GET",
			Body:   "{\"foo\": \"bar\"}",
		},
		Route{
			Route:  "/",
			Code:   200,
			Method: "GET",
			Body:   "{\"frank\": \"greco\"}",
		},
	})

	mockThree, _ := json.Marshal(mock{
		Route{
			Route:  "",
			Code:   200,
			Method: "GET",
			Body:   "{\"foo\": \"bar\"}",
		},
		Route{
			Route:  "foo",
			Code:   200,
			Method: "GET",
			Body:   "{\"foo\": \"bar\"}",
		},
		Route{
			Route:  "/foo/bar/car",
			Code:   200,
			Method: "POST",
			Body:   "{\"foo\": \"bar\"}",
		},
		Route{
			Route:  "/foo/bar/car",
			Code:   200,
			Method: "CUSTOM",
			Body:   "{\"foo\": \"bar\"}",
		},
	})

	mockFour, _ := json.Marshal(mock{
		Route{
			Route:  "/",
			Code:   200,
			Method: "GET",
			Body:   "{\"foo\": \"bar\"}",
		},
	})

	return []v1.ConfigMap{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "cm-one",
				Namespace: "foo",
			},
			Data: map[string]string{
				"response": string(mockOne),
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "cm-two",
				Namespace: "foo",
			},
			Data: map[string]string{
				"response": string(mockTwo),
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "cm-three",
				Namespace: "foo",
			},
			Data: map[string]string{
				"response": string(mockThree),
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "cm-four",
				Namespace: "bar",
			},
			Data: map[string]string{
				"response": string(mockOne),
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "cm-five",
				Namespace: "foo",
			},
			Data: map[string]string{
				"response": string(mockFour),
			},
		},
	}

}
