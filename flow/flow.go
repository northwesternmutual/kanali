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
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/northwesternmutual/kanali/controller"
	"github.com/opentracing/opentracing-go"
)

// Flow represents a series of steps
type Flow []Step

// Add takes a step and adds it to the flow
func (f *Flow) Add(steps ...Step) {
	for _, step := range steps {
		*f = append(*f, step)
	}
}

// Play iterates through every step that it has
// and exectutes them.
func (f *Flow) Play(ctx context.Context, ctlr *controller.Controller, w http.ResponseWriter, r *http.Request, resp *http.Response, trace opentracing.Span) error {

	logrus.Infof("flow is about to play")

	for _, step := range *f {
		logrus.Infof("playing step: %s", step.GetName())
		if err := step.Do(ctx, ctlr, w, r, resp, trace); err != nil {
			// An error was encourted during the flow's execution that Kanali is responsible for.
			// Because of this, we need to make sure we tag our opentracing span with that error
			trace.SetTag("error", true)
			trace.LogKV(
				"event", "error",
				"error.message", err.Error(),
			)
			return err
		}
	}
	return nil

}
