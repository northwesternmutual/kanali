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
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComputeTargetPath(t *testing.T) {
	tests := []struct {
		source, target, actual, expected string
	}{
		{
			source:   "/foo",
			target:   "/car",
			actual:   "/foo/bar",
			expected: "/car/bar",
		},
		{
			source:   "/",
			target:   "/car",
			actual:   "/",
			expected: "/car",
		},
		{
			source:   "/",
			target:   "/",
			actual:   "/",
			expected: "/",
		},
		{
			source:   "",
			target:   "",
			actual:   "",
			expected: "/",
		},
		{
			source:   "/foo/bar",
			target:   "/",
			actual:   "/foo/bar",
			expected: "/",
		},
		{
			source:   "/foo/bar",
			target:   "/foo",
			actual:   "/foo/bar",
			expected: "/foo",
		},
		{
			source:   "/foo",
			target:   "/",
			actual:   "/foo/bar",
			expected: "/bar",
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, ComputeTargetPath(test.source, test.target, test.actual))
	}
}

func BenchmarkComputeTargetPath(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ComputeTargetPath("/foo", "/car", "/foo/bar")
	}
}

func TestNormalizeURLPath(t *testing.T) {
	tests := []struct {
		path, expected string
	}{
		{
			path:     "",
			expected: "/",
		},
		{
			path:     "////",
			expected: "/",
		},
		{
			path:     "foo////bar",
			expected: "/foo/bar",
		},
		{
			path:     "foo",
			expected: "/foo",
		},
		{
			path:     "foo////",
			expected: "/foo",
		},
		{
			path:     "foo/",
			expected: "/foo",
		},
		{
			path:     "///foo////bar//",
			expected: "/foo/bar",
		},
		{
			path:     "/////https%3A%2F%2Fgoogle.com",
			expected: "/https%3A%2F%2Fgoogle.com",
		},
		{
			path:     "/foo",
			expected: "/foo",
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, NormalizeURLPath(test.path))
	}
}

func BenchmarkNormalizeURLPath(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NormalizeURLPath("///foo////bar//")
	}
}

func TestTransferResponse(t *testing.T) {
	headers := make(http.Header)
	headers.Add("foo", "bar")

	from := httptest.NewRecorder()
	from.Code = 200
	from.HeaderMap = headers
	_, err := from.Write([]byte("foo"))
	assert.Nil(t, err)

	to := httptest.NewRecorder()
	assert.Nil(t, TransferResponse(from, to))
	result := to.Result()
	body, _ := ioutil.ReadAll(result.Body)

	assert.Equal(t, from.Code, result.StatusCode)
	assert.Equal(t, from.HeaderMap, result.Header)
	assert.Equal(t, "foo", string(body))
}

func TestCopyHeader(t *testing.T) {
	original := http.Header(map[string][]string{
		"Foo": {"bar"},
	})
	assert.Equal(t, original, CloneHTTPHeader(original))
}
