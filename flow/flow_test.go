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
	"errors"
	"net/http"
	"testing"

	"github.com/northwesternmutual/kanali/logging"
	"github.com/northwesternmutual/kanali/metrics"
	"github.com/northwesternmutual/kanali/spec"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/stretchr/testify/assert"
)

func init() {
	logging.Init(nil)
}

type mockStep struct{}

func (s mockStep) GetName() string {
	return "mock step"
}

func (s mockStep) Do(ctx context.Context, proxy *spec.APIProxy, m *metrics.Metrics, w http.ResponseWriter, r *http.Request, resp *http.Response, trace opentracing.Span) error {
	return nil
}

type mockErrorStep struct{}

func (s mockErrorStep) GetName() string {
	return "mock error step"
}

func (s mockErrorStep) Do(ctx context.Context, proxy *spec.APIProxy, m *metrics.Metrics, w http.ResponseWriter, r *http.Request, resp *http.Response, trace opentracing.Span) error {
	return errors.New("forced error")
}

func TestAdd(t *testing.T) {
	f := &Flow{}
	f.Add(mockStep{}, mockStep{})
	f.Add(mockStep{})
	assert.Equal(t, len(*f), 3)
	for _, step := range *f {
		assert.Equal(t, step.GetName(), "mock step")
		assert.Nil(t, step.Do(context.Background(), nil, nil, nil, nil, nil, opentracing.StartSpan("test span")))
	}
}

func TestPlay(t *testing.T) {
	f := &Flow{}
	f.Add(mockStep{})
	assert.Nil(t, f.Play(context.Background(), nil, nil, nil, nil, nil, opentracing.StartSpan("test span")))
	f.Add(mockErrorStep{})
	assert.Error(t, f.Play(context.Background(), nil, nil, nil, nil, nil, opentracing.StartSpan("test span")))
}
