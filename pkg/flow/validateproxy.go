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
	"go.uber.org/zap"

	"github.com/northwesternmutual/kanali/pkg/errors"
	"github.com/northwesternmutual/kanali/pkg/logging"
	store "github.com/northwesternmutual/kanali/pkg/store/kanali/v2"
	"github.com/northwesternmutual/kanali/pkg/tags"
	"github.com/northwesternmutual/kanali/pkg/utils"
)

// ValidateProxyStep is factory that defines a step responsible for
// validating that an incoming request matches a proxy that Kanali
// has stored in memory
type validateProxyStep struct{}

func ValidateProxyStep() Step {
	return validateProxyStep{}
}

// GetName retruns the name of the ValidateProxyStep step
func (step validateProxyStep) Name() string {
	return "Validate Proxy"
}

// Do executes the logic of the ValidateProxyStep step
func (step validateProxyStep) Do(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	logger := logging.WithContext(r.Context())

	if proxy := store.ApiProxyStore().Get(utils.ComputeURLPath(r.URL)); proxy != nil {
		if span := opentracing.SpanFromContext(ctx); span != nil {
			span.SetTag(tags.KanaliProxyName, proxy.GetName())
			span.SetTag(tags.KanaliProxyNamespace, proxy.GetNamespace())
		}
		logger.Info("ApiProxy details",
			zap.String(tags.KanaliProxyName, proxy.GetName()),
			zap.String(tags.KanaliProxyNamespace, proxy.GetNamespace()),
		)
		return next()
	}

	logger.Warn(errors.ErrorProxyNotFound.Message)
	return errors.ErrorProxyNotFound
}
