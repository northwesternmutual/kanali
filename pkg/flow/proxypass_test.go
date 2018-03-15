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

package flow

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/labels"

	"github.com/northwesternmutual/kanali/pkg/apis/kanali.io/v2"
	"github.com/northwesternmutual/kanali/test/builder"
)

func TestGetServiceLabelSet(t *testing.T) {
	tests := []struct {
		proxy    *v2.ApiProxy
		headers  http.Header
		labels   labels.Set
		defaults map[string]string
	}{
		{
			proxy: builder.NewApiProxy("foo", "bar").WithTargetBackendDynamicService(8080,
				v2.Label{Name: "foo", Value: "bar"},
			).NewOrDie(),
			headers: nil,
			labels: map[string]string{
				"foo": "bar",
			},
			defaults: make(map[string]string),
		},
		{
			proxy:    builder.NewApiProxy("foo", "bar").WithTargetBackendDynamicService(8080).NewOrDie(),
			headers:  nil,
			labels:   map[string]string{},
			defaults: make(map[string]string),
		},
		{
			proxy: builder.NewApiProxy("foo", "bar").WithTargetBackendDynamicService(8080,
				v2.Label{Name: "foo", Value: "bar"},
				v2.Label{Name: "bar", Header: "foo"},
			).NewOrDie(),
			headers: make(http.Header),
			labels: map[string]string{
				"foo": "bar",
				"bar": "",
			},
			defaults: make(map[string]string),
		},
		{
			proxy: builder.NewApiProxy("foo", "bar").WithTargetBackendDynamicService(8080,
				v2.Label{Name: "foo", Value: "bar"},
				v2.Label{Name: "bar", Header: "foo"},
			).NewOrDie(),
			headers: http.Header(map[string][]string{
				"Foo": {"car"},
			}),
			labels: map[string]string{
				"foo": "bar",
				"bar": "car",
			},
			defaults: make(map[string]string),
		},
		{
			proxy: builder.NewApiProxy("foo", "bar").WithTargetBackendDynamicService(8080,
				v2.Label{Name: "bar", Header: "foo"},
			).NewOrDie(),
			headers: http.Header(map[string][]string{
				"Foo": {"car"},
			}),
			labels: map[string]string{
				"bar": "car",
			},
			defaults: make(map[string]string),
		},
		{
			proxy: builder.NewApiProxy("foo", "bar").WithTargetBackendDynamicService(8080,
				v2.Label{Name: "foo", Value: "bar"},
				v2.Label{Name: "bar", Header: "foo"},
			).NewOrDie(),
			headers: make(http.Header),
			labels: map[string]string{
				"foo": "bar",
				"bar": "default",
			},
			defaults: map[string]string{
				"foo": "default",
			},
		},
		{
			proxy: builder.NewApiProxy("foo", "bar").WithTargetBackendDynamicService(8080,
				v2.Label{Name: "foo", Value: "bar"},
				v2.Label{Name: "bar", Header: "foo"},
			).NewOrDie(),
			headers: http.Header(map[string][]string{
				"Foo": {"car"},
			}),
			labels: map[string]string{
				"foo": "bar",
				"bar": "car",
			},
			defaults: map[string]string{
				"foo": "default",
			},
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.labels, getServiceLabelSet(test.proxy, test.headers, test.defaults))
	}
}

func BenchmarkGetServiceLabelSet(b *testing.B) {
	proxy := builder.NewApiProxy("foo", "bar").WithTargetBackendDynamicService(8080,
		v2.Label{Name: "foo", Value: "bar"},
		v2.Label{Name: "bar", Header: "foo"},
	).NewOrDie()
	headers := make(http.Header)
	defaults := map[string]string{
		"foo": "default",
	}

	for n := 0; n < b.N; n++ {
		getServiceLabelSet(proxy, headers, defaults)
	}
}

func TestCopyBuffer(t *testing.T) {

}

func TestCopyHeader(t *testing.T) {
  original := http.Header(map[string][]string{
    "Foo": []string{"bar"},
  })
  copy := make(http.Header)
  copyHeader(copy, original)
  assert.Equal(t, 1, len(copy))
  assert.Equal(t, "bar", copy["Foo"][0])
  delete(original, "Foo")
  assert.Equal(t, 0, len(original))
  assert.Equal(t, 0, len(copy))
}
