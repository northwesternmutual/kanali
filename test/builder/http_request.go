// Copyright (c) 2018 Northwestern Mutual.
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

package builder

import (
	"net/http"
	"net/url"
)

type HTTPRequestBuilder struct {
	curr http.Request
}

func NewHTTPRequest() *HTTPRequestBuilder {
	return &HTTPRequestBuilder{
		curr: http.Request{
			Header: make(http.Header),
			URL:    &url.URL{},
		},
	}
}

func (b *HTTPRequestBuilder) WithMethod(method string) *HTTPRequestBuilder {
	b.curr.Method = method
	return b
}

func (b *HTTPRequestBuilder) WithHeader(key string, values ...string) *HTTPRequestBuilder {
	if b.curr.Header == nil {
		b.curr.Header = map[string][]string{}
	}
	for _, val := range values {
		b.curr.Header.Add(key, val)
	}
	return b
}

func (b *HTTPRequestBuilder) WithHost(host string) *HTTPRequestBuilder {
	u, _ := url.Parse(host)
	b.curr.URL.Scheme = u.Scheme
	b.curr.URL.Host = u.Host
	return b
}

func (b *HTTPRequestBuilder) WithPath(path string) *HTTPRequestBuilder {
	u, _ := url.ParseRequestURI(path)
	b.curr.URL.Path = u.Path
	b.curr.URL.RawPath = u.EscapedPath()
	return b
}

func (b *HTTPRequestBuilder) NewOrDie() *http.Request {
	return &b.curr
}
