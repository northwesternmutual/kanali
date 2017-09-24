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

package steps

import (
	"context"
	"errors"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/northwesternmutual/kanali/metrics"
	"github.com/northwesternmutual/kanali/spec"
	"github.com/northwesternmutual/kanali/tracer"
	"github.com/northwesternmutual/kanali/utils"
	"github.com/opentracing/opentracing-go"
)

// ValidateProxyStep is factory that defines a step responsible for
// validating that an incoming request matches a proxy that Kanali
// has stored in memory
type ValidateProxyStep struct{}

// GetName retruns the name of the ValidateProxyStep step
func (step ValidateProxyStep) GetName() string {
	return "Validate Proxy"
}

// Do executes the logic of the ValidateProxyStep step
func (step ValidateProxyStep) Do(ctx context.Context, proxy *spec.APIProxy, m *metrics.Metrics, w http.ResponseWriter, r *http.Request, resp *http.Response, trace opentracing.Span) error {

	untypedProxy, err := spec.ProxyStore.Get(r.URL.EscapedPath())
	if err != nil || untypedProxy == nil {
		if err != nil {
			logrus.Error(err.Error())
		}

		trace.SetTag(tracer.KanaliProxyName, "unknown")
		trace.SetTag(tracer.KanaliProxyNamespace, "unknown")

		m.Add(
			metrics.Metric{Name: "proxy_name", Value: "unknown", Index: true},
			metrics.Metric{Name: "proxy_namespace", Value: "unknown", Index: true},
		)

		return utils.StatusError{Code: http.StatusNotFound, Err: errors.New("proxy not found")}
	}

	typedProxy, _ := untypedProxy.(spec.APIProxy)

	proxy = &typedProxy

	trace.SetTag(tracer.KanaliProxyName, proxy.ObjectMeta.Name)
	trace.SetTag(tracer.KanaliProxyNamespace, proxy.ObjectMeta.Namespace)

	m.Add(
		metrics.Metric{Name: "proxy_name", Value: proxy.ObjectMeta.Name, Index: true},
		metrics.Metric{Name: "proxy_namespace", Value: proxy.ObjectMeta.Namespace, Index: true},
	)

	return nil

}
