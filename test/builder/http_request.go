package builder

import (
	"net"
	"net/http"
	"net/url"
)

type HTTPRequestBuilder struct {
	curr http.Request
}

func NewHTTPRequest(method string) *HTTPRequestBuilder {
	return &HTTPRequestBuilder{
		curr: http.Request{
			Method: method,
			URL:    &url.URL{},
		},
	}
}

func (b *HTTPRequestBuilder) WithHTTP() *HTTPRequestBuilder {
	b.curr.URL.Scheme = "http"
	return b
}

func (b *HTTPRequestBuilder) WithHostPort(host, port string) *HTTPRequestBuilder {
	b.curr.URL.Host = net.JoinHostPort(host, port)
	return b
}

func (b *HTTPRequestBuilder) WithHTTPS() *HTTPRequestBuilder {
	b.curr.URL.Scheme = "https"
	return b
}

func (b *HTTPRequestBuilder) WithPath(path string) *HTTPRequestBuilder {
	b.curr.URL.Path = path
	return b
}

func (b *HTTPRequestBuilder) NewOrDie() *http.Request {
	return &b.curr
}
