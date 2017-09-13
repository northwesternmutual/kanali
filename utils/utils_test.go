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

package utils

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/kubernetes/pkg/api"
)

func TestComputeTargetPath(t *testing.T) {

	assert := assert.New(t)

	assert.Equal("/", ComputeTargetPath("/foo/bar", "", "/foo/bar"), "path not what expected")
	assert.Equal("/", ComputeTargetPath("/foo/bar", "/", "/foo/bar"), "path not what expected")
	assert.Equal("/foo", ComputeTargetPath("/foo/bar", "/foo", "/foo/bar"), "path not what expected")
	assert.Equal("/foo/bar", ComputeTargetPath("/foo/bar", "/foo", "/foo/bar/bar"), "path not what expected")
	assert.Equal("/bar", ComputeTargetPath("/foo/bar", "", "/foo/bar/bar"), "path not what expected")
	assert.Equal("/accounts", ComputeTargetPath("/api/v1/example-two", "/", "/api/v1/example-two/accounts"), "path not what expected")
	assert.Equal("/accounts", ComputeTargetPath("/api/v1/example-two", "/", "/api/v1/example-two/accounts/"), "path not what expected")
	assert.Equal("/accounts", ComputeTargetPath("/api/v1/example-two", "", "/api/v1/example-two/accounts/"), "path not what expected")
	assert.Equal("/accounts", ComputeTargetPath("/api/v1/example-two/", "/", "/api/v1/example-two/accounts/"), "path not what expected")
	assert.Equal("/accounts", ComputeTargetPath("/api/v1/example-two/", "", "/api/v1/example-two/accounts/"), "path not what expected")
	assert.Equal("/accounts", ComputeTargetPath("/api/v1/example-two/", "", "/api/v1/example-two/accounts"), "path not what expected")
	assert.Equal("/", ComputeTargetPath("/", "", "/"), "path not what expected")
	assert.Equal("/", ComputeTargetPath("/", "/", "/"), "path not what expected")

}

func TestAbsPath(t *testing.T) {
	p, _ := GetAbsPath("/")
	assert.Equal(t, "", p)
	p, _ = GetAbsPath("/foo/")
	assert.Equal(t, "/foo", p)
	p, _ = GetAbsPath("//")
	assert.Equal(t, "", p)
}

func TestIsValidHTTPMethod(t *testing.T) {
	assert.False(t, IsValidHTTPMethod("foo"))
	assert.True(t, IsValidHTTPMethod("GET"))
	assert.True(t, IsValidHTTPMethod("get"))
	assert.True(t, IsValidHTTPMethod("POST"))
	assert.True(t, IsValidHTTPMethod("post"))
}

func TestFlattenHTTPHeaders(t *testing.T) {
	h := http.Header{
		"Foo": []string{"bar"},
		"Bar": []string{"foo"},
	}
	assert.Equal(t, FlattenHTTPHeaders(h), map[string]string{
		"Foo": "bar",
		"Bar": "foo",
	})
	assert.Nil(t, FlattenHTTPHeaders(nil))
}

func TestDupReaderAndString(t *testing.T) {
	closer := ioutil.NopCloser(bytes.NewReader([]byte("test string")))
	closerTwo, str, _ := DupReaderAndString(closer)
	assert.Equal(t, str, "test string")
	data, _ := ioutil.ReadAll(closerTwo)
	assert.Equal(t, string(data), "test string")
}

func TestOmitHeaderValues(t *testing.T) {

	assert := assert.New(t)

	h := http.Header{
		"One":   []string{"two"},
		"Three": []string{"four"},
	}

	copy := OmitHeaderValues(h, "ommitted", "one")
	assert.Equal(h, http.Header{
		"One":   []string{"two"},
		"Three": []string{"four"},
	}, "original map should not change")
	assert.Equal(copy, http.Header{
		"One":   []string{"ommitted"},
		"Three": []string{"four"},
	}, "map should be equal")

	copy = OmitHeaderValues(h, "ommitted", "one", "foo", "bar")
	assert.Equal(copy, http.Header{
		"One":   []string{"ommitted"},
		"Three": []string{"four"},
	}, "map should be equal")

	copy = OmitHeaderValues(h, "ommitted")
	assert.Equal(copy, http.Header{
		"One":   []string{"two"},
		"Three": []string{"four"},
	}, "original map should not change")

	copy = OmitHeaderValues(nil, "ommitted")
	assert.Nil(copy, "map should be equal")

}

func TestCompareObjectMeta(t *testing.T) {
	c1 := api.ObjectMeta{
		Name:      "foo",
		Namespace: "bar",
	}
	c2 := api.ObjectMeta{
		Name:      "bar",
		Namespace: "foo",
	}
	c3 := api.ObjectMeta{
		Name:      "foo",
		Namespace: "car",
	}
	c4 := api.ObjectMeta{
		Name:      "bar",
		Namespace: "car",
	}

	assert.True(t, CompareObjectMeta(c1, c1))
	assert.False(t, CompareObjectMeta(c1, c2))
	assert.False(t, CompareObjectMeta(c1, c3))
	assert.False(t, CompareObjectMeta(c3, c4))
}

func TestNormalizePath(t *testing.T) {
	assert.Equal(t, "/foo/bar", NormalizePath("foo////bar"))
	assert.Equal(t, "/foo", NormalizePath("foo"))
	assert.Equal(t, "/foo", NormalizePath("foo////"))
	assert.Equal(t, "/foo/bar", NormalizePath("///foo////bar//"))
	assert.Equal(t, "/", NormalizePath(""))
	assert.Equal(t, "/", NormalizePath("////"))
}
