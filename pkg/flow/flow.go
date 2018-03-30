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

	opentracing "github.com/opentracing/opentracing-go"
	tracelog "github.com/opentracing/opentracing-go/log"
	"go.uber.org/zap"

	"github.com/northwesternmutual/kanali/pkg/log"
	"github.com/northwesternmutual/kanali/pkg/tags"
)

// Flow represents a list of Steps.
type Flow []Step

// New will create a new Flow and return its reference.
func New() *Flow {
	return &Flow{}
}

// Add will add an arbitary number of Steps to the Flow
// in the order in which they were passed to this method.
func (f *Flow) Add(steps ...Step) *Flow {
	for _, s := range steps {
		*f = append(*f, s)
	}
	return f
}

// Play will execute each Step in the Flow in sequential order.
// If the execution of any Step returns an error, this error will
// be returned immedietaly and the execution of any latter Step
// will not occur.
func (f *Flow) Play(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	logger := log.WithContext(r.Context())

	for _, step := range *f {
		logger.With(
			zap.String("step.name", step.Name()),
		).Debug("playing step")
		err := step.Do(ctx, w, r)
		if err == nil {
			continue
		}
		if span := opentracing.SpanFromContext(ctx); span != nil {
			span.SetTag(tags.Error, true)
			span.LogFields(
				tracelog.String("event", tags.Error),
				tracelog.String("error.message", err.Error()),
			)
		}
		return err
	}

	return nil
}
