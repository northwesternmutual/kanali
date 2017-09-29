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

package handlers

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalize(t *testing.T) {
	r1 := &http.Request{
		URL: &url.URL{
			Path: "///foo//bar/car",
		},
	}
	r2 := &http.Request{
		URL: &url.URL{
			Path: "foo//bar/car/",
		},
	}
	r3 := &http.Request{
		URL: &url.URL{
			Path: "",
		},
	}
	r4 := &http.Request{
		URL: &url.URL{
			Path: "////",
		},
	}
	normalize(r1)
	normalize(r2)
	normalize(r3)
	normalize(r4)

	assert.Equal(t, "/foo/bar/car", r1.URL.Path)
	assert.Equal(t, "/foo/bar/car", r2.URL.Path)
	assert.Equal(t, "/", r3.URL.Path)
	assert.Equal(t, "/", r4.URL.Path)

}