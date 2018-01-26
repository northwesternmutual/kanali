package flow

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/northwesternmutual/kanali/pkg/apis/kanali.io/v2"
	"github.com/northwesternmutual/kanali/pkg/errors"
)

func TestEndpointDetails(t *testing.T) {
	p1 := &v2.ApiProxy{Spec: v2.ApiProxySpec{Target: v2.Target{Backend: v2.Backend{Endpoint: "https://foo.bar.com/baz"}}}}
	p2 := &v2.ApiProxy{Spec: v2.ApiProxySpec{Target: v2.Target{Backend: v2.Backend{Endpoint: "http://foo.bar.com/baz"}}}}
	p3 := &v2.ApiProxy{Spec: v2.ApiProxySpec{Target: v2.Target{Backend: v2.Backend{Endpoint: "foo.bar.com/baz"}}}}

	s1, h1, e1 := proxyPassStep{proxy: p1}.endpointDetails()
	assert.Nil(t, e1)
	assert.Equal(t, "https", s1)
	assert.Equal(t, "foo.bar.com", h1)

	s2, h2, e2 := proxyPassStep{proxy: p2}.endpointDetails()
	assert.Nil(t, e2)
	assert.Equal(t, "http", s2)
	assert.Equal(t, "foo.bar.com", h2)

	_, _, e3 := proxyPassStep{proxy: p3}.endpointDetails()
	assert.Equal(t, e3, errors.ErrorApiProxyBackendEndpointMalformed)
}

func BenchmarkEndpointDetails(b *testing.B) {
	for n := 0; n < b.N; n++ {
		proxyPassStep{
			proxy: &v2.ApiProxy{Spec: v2.ApiProxySpec{Target: v2.Target{Backend: v2.Backend{Endpoint: "https://foo.bar.com/baz"}}}},
		}.endpointDetails()
	}
}
