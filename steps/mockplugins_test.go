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

// package
package steps

import (
	"context"
	"errors"
	"net/http"

	"github.com/northwesternmutual/kanali/metrics"
	"github.com/northwesternmutual/kanali/spec"
	"github.com/opentracing/opentracing-go"
)

type fakePanicPlugin struct{}

func (plugin fakePanicPlugin) OnRequest(ctx context.Context, m *metrics.Metrics, p spec.APIProxy, r *http.Request, span opentracing.Span) error {
	panic("intentional")
}

func (plugin fakePanicPlugin) OnResponse(ctx context.Context, m *metrics.Metrics, p spec.APIProxy, r *http.Request, resp *http.Response, span opentracing.Span) error {
	panic("intentional")
}

type fakeSuccessPlugin struct{}

func (plugin fakeSuccessPlugin) OnRequest(ctx context.Context, m *metrics.Metrics, p spec.APIProxy, r *http.Request, span opentracing.Span) error {
	return nil
}

func (plugin fakeSuccessPlugin) OnResponse(ctx context.Context, m *metrics.Metrics, p spec.APIProxy, r *http.Request, resp *http.Response, span opentracing.Span) error {
	return nil
}

type fakeErrorPlugin struct{}

func (plugin fakeErrorPlugin) OnRequest(ctx context.Context, m *metrics.Metrics, p spec.APIProxy, r *http.Request, span opentracing.Span) error {
	return errors.New("error")
}

func (plugin fakeErrorPlugin) OnResponse(ctx context.Context, m *metrics.Metrics, p spec.APIProxy, r *http.Request, resp *http.Response, span opentracing.Span) error {
	return errors.New("error")
}
