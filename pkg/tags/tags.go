package pkg

const (
	// Error is the opentracing tag name that represents an error
	Error = "error"

	// KanaliProxyName is the opentracing tag name that represents an APIProxy name
	KanaliProxyName = "kanali.proxy.name"
	// KanaliProxyNamespace is the opentracing tag name that represents an APIProxy namespace
	KanaliProxyNamespace = "kanali.proxy.namespace"

	// HTTPRequest is the opentracing tag name that represents the existence on an HTTP request
	HTTPRequest = "http.request"

	// HTTPRequestRemoteAddress is the ip address of the caller
	HTTPRequestRemoteAddress = "http.request.remote_address"
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
