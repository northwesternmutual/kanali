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

	"github.com/northwesternmutual/kanali/pkg/apis/kanali.io/v2"
	"github.com/northwesternmutual/kanali/pkg/logging"
	"github.com/northwesternmutual/kanali/pkg/metrics"
	"github.com/northwesternmutual/kanali/pkg/steps"
	"github.com/northwesternmutual/kanali/pkg/tags"
	opentracing "github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"k8s.io/client-go/informers/core"
)

type flow []steps.Step

func New() *flow {
	return &flow{}
}

func (f *flow) Add(steps ...steps.Step) {
	for _, s := range steps {
		*f = append(*f, s)
	}
}

func (f *flow) Play(ctx context.Context, proxy *v2.ApiProxy, k8sCoreClient core.Interface, metrics *metrics.Metrics, w http.ResponseWriter, r *http.Request, resp *http.Response, trace opentracing.Span) error {
	logger := logging.WithContext(ctx)

	for _, step := range *f {
		logger.With(
			zap.String("step.name", step.GetName()),
		).Debug("playing step")
		err := step.Do(ctx, proxy, k8sCoreClient, metrics, w, r, resp, trace)
		if err == nil {
			continue
		}
		trace.SetTag(tags.Error, true)
		trace.LogKV(
			"event", tags.Error,
			"error.message", err.Error(),
		)
		return err
	}
	return nil
}
