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

	"github.com/Sirupsen/logrus"
	"github.com/northwesternmutual/kanali/metrics"
	"github.com/northwesternmutual/kanali/plugins"
	"github.com/northwesternmutual/kanali/spec"
	"github.com/opentracing/opentracing-go"
)

// PluginsOnRequestStep is factory that defines a step responsible for
// executing the on request lifecycle hook for all the defined plugins
type PluginsOnRequestStep struct{}

// GetName retruns the name of the PluginsOnRequestStep step
func (step PluginsOnRequestStep) GetName() string {
	return "Plugin OnRequest"
}

// Do executes the logic of the PluginsOnRequestStep step
func (step PluginsOnRequestStep) Do(ctx context.Context, proxy *spec.APIProxy, m *metrics.Metrics, w http.ResponseWriter, r *http.Request, resp *http.Response, trace opentracing.Span) error {

	for _, plugin := range proxy.Spec.Plugins {
		p, err := plugins.GetPlugin(plugin)
		if err != nil {
			return err
		}
		if err := doOnRequest(ctx, m, plugin.Name, *proxy, r, trace, *p); err != nil {
			return err
		}
	}
	return nil
}

func doOnRequest(ctx context.Context, m *metrics.Metrics, name string, proxy spec.APIProxy, req *http.Request, span opentracing.Span, p plugins.Plugin) (e error) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorf("OnRequest paniced: %v", r)
			e = errors.New("OnRequest paniced")
		}
	}()

	sp := opentracing.StartSpan(fmt.Sprintf("PLUGIN: ON_REQUEST: %s", name), opentracing.ChildOf(span.Context()))
	defer sp.Finish()

	return p.OnRequest(ctx, m, proxy, req, sp)
}
