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

package flow

import (
	"context"
	"fmt"
	"net/http"

	"github.com/northwesternmutual/kanali/logging"
	"github.com/northwesternmutual/kanali/metrics"
	"github.com/northwesternmutual/kanali/spec"
	"github.com/northwesternmutual/kanali/tracer"
	"github.com/opentracing/opentracing-go"
)

// Flow is a list of steps
type Flow []Step

// Add appends a step to a flow
func (f *Flow) Add(steps ...Step) {
	for _, step := range steps {
		*f = append(*f, step)
	}
}

// Play executes all step in a flow in the order they were added.
func (f *Flow) Play(ctx context.Context, proxy *spec.APIProxy, metrics *metrics.Metrics, w http.ResponseWriter, r *http.Request, resp *http.Response, trace opentracing.Span) error {
	logger := logging.WithContext(ctx)

	logger.Debug("flow is about to play")
	for _, step := range *f {
		logger.Debug(fmt.Sprintf("playing step %s", step.GetName()))
		if err := step.Do(ctx, proxy, metrics, w, r, resp, trace); err != nil {
			trace.SetTag(tracer.Error, true)
			trace.LogKV(
				"event", "error",
				"error.message", err.Error(),
			)
			return err
		}
	}
	return nil
}
