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
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/northwesternmutual/kanali/pkg/utils/pool"
)

func Write(w io.Writer, msg string) error {
	_, err := w.Write([]byte(msg + "\n"))
	return err
}

// ComputeTargetPath calculates the path to be used when
// proxying to an upstream service.
func ComputeTargetPath(source, target, actual string) string {
	w := pool.GetBuffer()
	defer pool.PutBuffer(w)

	return doComputeTargetPath(w, source, target, actual)
}

func doComputeTargetPath(b *bytes.Buffer, source, target, actual string) string {
	var i int

	for i < len(source) && actual[i] == source[i] {
		i++
	}

	if len(target) != 1 || target[0] != '/' {
		for j := 0; j < len(target); j++ {
			b.WriteByte(target[j])
		}
	}

	for i < len(actual) {
		b.WriteByte(actual[i])
		i++
	}

	if b.Len() == 0 {
		b.WriteByte('/')
	}

	return b.String()
}

// ComputeURLPath will correct a URL path that might be valid but not ideally formatted
func ComputeURLPath(u *url.URL) string {
	return NormalizeURLPath(u.EscapedPath())
}

// NormalizeURLPath normalizes a path in the following ways:
// 1. remove duplicate slashes
// 2. remove trailing slash
// 3. ensure leading slash
func NormalizeURLPath(path string) string {
	w := pool.GetBuffer()
	defer pool.PutBuffer(w)

	return doNormalizeURLPath(w, path)
}

func doNormalizeURLPath(b *bytes.Buffer, path string) string {
	var i int

	b.WriteByte('/')

	if len(path) > 0 && path[0] == '/' {
		i++
	}

	for i < len(path) {
		if path[i] != '/' || path[i-1] != '/' {
			b.WriteByte(path[i])
		}
		i++
	}

	result := b.String()
	if len(result) > 1 && result[len(result)-1] == '/' {
		return result[:len(result)-1]
	}
	return result
}

func TransferResponse(from *httptest.ResponseRecorder, to http.ResponseWriter) error {
	for k, v := range from.Header() {
		to.Header()[k] = v
	}
	to.WriteHeader(from.Code)
	_, err := to.Write(from.Body.Bytes())
	return err
}

func CloneHTTPHeader(h http.Header) http.Header {
	h2 := make(http.Header, len(h))
	for k, vv := range h {
		vv2 := make([]string, len(vv))
		copy(vv2, vv)
		h2[k] = vv2
	}
	return h2
}
