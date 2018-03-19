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

package middleware

import (
	"context"
	"net/http"

	opentracing "github.com/opentracing/opentracing-go"

	"github.com/northwesternmutual/kanali/pkg/errors"
	"github.com/northwesternmutual/kanali/pkg/flow"
	"github.com/northwesternmutual/kanali/pkg/log"
	"github.com/northwesternmutual/kanali/pkg/tags"
	"github.com/northwesternmutual/kanali/pkg/tracer"
)

// Gateway is an http.HandlerFunc that preforms the core functions of the Kanali gateway.
// This is meant to be the final http.Handler if multiple middlewares are to be used.
func Gateway(w http.ResponseWriter, r *http.Request) {
	logger := log.WithContext(r.Context())
	span := tracer.StartSpan(r)
	ctx := opentracing.ContextWithSpan(context.Background(), span)
	defer span.Finish()

	tracer.HydrateSpanFromRequest(r, span)

	if err := flow.New().Add(
		flow.ValidateProxyStep(),
		flow.PluginsOnRequestStep(),
		flow.MockTargetStep(),
		flow.ProxyPassStep(),
		flow.PluginsOnResponseStep(),
	).Play(ctx, w, r); err != nil {
		err, data := errors.ToJSON(err)
		span.SetTag(tags.HTTPResponseBody, string(data))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(err.Status)
		if _, err := w.Write(data); err != nil {
			logger.Error(err.Error())
		}
	}

	// TODO:
	//tracer.HydrateSpanFromResponse(r, span)
}
