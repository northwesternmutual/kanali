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
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"strings"
)

// ComputeTargetPath calcuates the target or destination path based on the incoming path,
// desired target path prefix and the assicated proxy
//
//
func ComputeTargetPath(proxyPath, proxyTarget, requestPath string) string {
	var buffer bytes.Buffer

	if len(strings.SplitAfter(requestPath, proxyPath)) == 0 {
		buffer.WriteString("/")
	} else if proxyTarget != "/" {
		buffer.WriteString(proxyTarget)
	}

	buffer.WriteString(strings.SplitAfter(requestPath, proxyPath)[1])

	if len(buffer.Bytes()) == 0 {
		return "/"
	}

	return buffer.String()
}

// ComputeURLPath will correct a URL path that might be valid but not ideally formatted
func ComputeURLPath(u *url.URL) string {
	return NormalizeURLPath(u.EscapedPath())
}

func NormalizeURLPath(path string) string {
	if len(path) < 1 {
		return "/"
	}

	path = regexp.MustCompile(`/{2,}`).ReplaceAllString(path, "/")

	if strings.HasSuffix(path, "/") {
		path = path[:len(path)-1]
	}

	if len(path) < 1 {
		return "/"
	}

	if path[0] != '/' {
		path = "/" + path
	}

	return path
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
