package builder

import (
	"net/http"
	"net/http/httptest"
)

type HTTPServerBuilder struct {
	curr httptest.Server
}

func NewHTTPServer(handler http.Handler) *HTTPServerBuilder {
	return &HTTPServerBuilder{
		curr: *httptest.NewUnstartedServer(handler),
	}
}

func (b *HTTPServerBuilder) WithTLS(ca []byte) *HTTPServerBuilder {
	return b
}

func (b *HTTPServerBuilder) WithCA(cert, key []byte) *HTTPServerBuilder {
	return b
}

func (b *HTTPServerBuilder) NewOrDie() *httptest.Server {
	return &b.curr
}
