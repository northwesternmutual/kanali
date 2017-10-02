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
	"fmt"
	"net/http"

	"github.com/northwesternmutual/kanali/logging"
	"github.com/northwesternmutual/kanali/metrics"
	"github.com/northwesternmutual/kanali/plugins"
	"github.com/northwesternmutual/kanali/spec"
	"github.com/opentracing/opentracing-go"
)

// PluginsOnResponseStep is factory that defines a step responsible for
// executing the on response lifecycle hook for all the defined plugins
type PluginsOnResponseStep struct{}

// GetName retruns the name of the PluginsOnResponseStep step
func (step PluginsOnResponseStep) GetName() string {
	return "Plugin OnResponse"
}

// Do executes the logic of the PluginsOnResponseStep step
func (step PluginsOnResponseStep) Do(ctx context.Context, proxy *spec.APIProxy, m *metrics.Metrics, w http.ResponseWriter, r *http.Request, resp *http.Response, trace opentracing.Span) error {

	for _, plugin := range proxy.Spec.Plugins {
		p, err := plugins.GetPlugin(plugin)
		if err != nil {
			return err
		}
		if err := doOnResponse(ctx, m, plugin.Name, *proxy, r, resp, trace, *p); err != nil {
			return err
		}
	}

	return nil
}

func doOnResponse(ctx context.Context, m *metrics.Metrics, name string, proxy spec.APIProxy, req *http.Request, resp *http.Response, span opentracing.Span, p plugins.Plugin) (e error) {
	logger := logging.WithContext(ctx)

	defer func() {
		if r := recover(); r != nil {
			logger.Error(fmt.Sprintf("OnResponse paniced: %v", r))
			e = errors.New("OnResponse paniced")
		}
	}()

	sp := opentracing.StartSpan(fmt.Sprintf("PLUGIN: ON_RESPONSE: %s", name), opentracing.ChildOf(span.Context()))
	defer sp.Finish()

	return p.OnResponse(ctx, m, proxy, req, resp, sp)
}
