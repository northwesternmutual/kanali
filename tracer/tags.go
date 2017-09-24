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

package tracer

const (
	// Error is the opentracing tag name that represents an error
	Error = "error"

	// KanaliProxyName is the opentracing tag name that represents an APIProxy name
	KanaliProxyName = "kanali.proxy.name"
	// KanaliProxyNamespace is the opentracing tag name that represents an APIProxy namespace
	KanaliProxyNamespace = "kanali.proxy.namespace"

	// HTTPRequest is the opentracing tag name that represents the existence on an HTTP request
	HTTPRequest = "http.request"

	// HTTPRequestMethod is the opentracing tag name that represents an HTTP request method
	HTTPRequestMethod = "http.request.method"
	// HTTPRequestBody is the opentracing tag name that represents an HTTP request body
	HTTPRequestBody = "http.request.body"
	// HTTPRequestHeaders is the opentracing tag name that represents an HTTP request headers
	HTTPRequestHeaders = "http.request.headers"

	// HTTPRequestURLPath is the opentracing tag name that represents an HTTP URL path
	HTTPRequestURLPath = "http.request.url.path"
	// HTTPRequestURLHost is the opentracing tag name that represents an HTTP URL host
	HTTPRequestURLHost = "http.request.url.host"
	// HTTPRequestURLQuery is the opentracing tag name that represents an HTTP URL query
	HTTPRequestURLQuery = "http.request.url.query"

	// HTTPResponse is the opentracing tag name that represents the existence on an HTTP response
	HTTPResponse = "http.response"

	// HTTPResponseBody is the opentracing tag name that represents an HTTP response body
	HTTPResponseBody = "http.response.body"
	// HTTPResponseHeaders is the opentracing tag name that represents an HTTP response headers
	HTTPResponseHeaders = "http.response.headers"

	// HTTPResponseStatusCode is the opentracing tag name that represents an HTTP response status code
	HTTPResponseStatusCode = "http.response.status.code"
)
