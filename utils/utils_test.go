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
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComputeTargetPath(t *testing.T) {

	assert := assert.New(t)

	assert.Equal("/", ComputeTargetPath("/foo/bar", "", "/foo/bar"))
	assert.Equal("/", ComputeTargetPath("/foo/bar", "/", "/foo/bar"))
	assert.Equal("/foo", ComputeTargetPath("/foo/bar", "/foo", "/foo/bar"))
	assert.Equal("/foo/bar", ComputeTargetPath("/foo/bar", "/foo", "/foo/bar/bar"))
	assert.Equal("/bar", ComputeTargetPath("/foo/bar", "", "/foo/bar/bar"))
	assert.Equal("/accounts", ComputeTargetPath("/api/v1/example-two", "/", "/api/v1/example-two/accounts"))
	assert.Equal("/accounts", ComputeTargetPath("/api/v1/example-two", "/", "/api/v1/example-two/accounts/"))
	assert.Equal("/accounts", ComputeTargetPath("/api/v1/example-two", "", "/api/v1/example-two/accounts/"))
	assert.Equal("/accounts", ComputeTargetPath("/api/v1/example-two/", "/", "/api/v1/example-two/accounts/"))
	assert.Equal("/accounts", ComputeTargetPath("/api/v1/example-two/", "", "/api/v1/example-two/accounts/"))
	assert.Equal("/accounts", ComputeTargetPath("/api/v1/example-two/", "", "/api/v1/example-two/accounts"))
	assert.Equal("/", ComputeTargetPath("/", "", "/"), "path not what expected")
	assert.Equal("/", ComputeTargetPath("/", "/", "/"), "path not what expected")

}

func TestAbsPath(t *testing.T) {

	assert := assert.New(t)

	p, _ := GetAbsPath("/")
	assert.Equal("", p)

	p, _ = GetAbsPath("/foo/")
	assert.Equal("/foo", p)

	p, _ = GetAbsPath("//")
	assert.Equal("", p)

}

func TestOmitHeaderValues(t *testing.T) {

	assert := assert.New(t)

	h := http.Header{
		"One":   []string{"two"},
		"Three": []string{"four"},
	}

	copy := OmitHeaderValues(h, "omitted", "one")
	assert.Equal(h, http.Header{
		"One":   []string{"two"},
		"Three": []string{"four"},
	}, "original map should not change")
	assert.Equal(copy, http.Header{
		"One":   []string{"omitted"},
		"Three": []string{"four"},
	}, "map should be equal")

	copy = OmitHeaderValues(h, "omitted", "one", "foo", "bar")
	assert.Equal(copy, http.Header{
		"One":   []string{"omitted"},
		"Three": []string{"four"},
	}, "map should be equal")

	copy = OmitHeaderValues(h, "omitted")
	assert.Equal(copy, http.Header{
		"One":   []string{"two"},
		"Three": []string{"four"},
	}, "original map should not change")

	copy = OmitHeaderValues(nil, "omitted")
	assert.Nil(copy, "map should be equal")

}
