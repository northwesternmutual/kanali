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

package flow

import (
	"context"
	"errors"
	"net/http"
	"testing"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/mocktracer"
	"github.com/stretchr/testify/assert"
)

type mockStep struct{}

func (s mockStep) Name() string {
	return "mock step"
}

func (s mockStep) Do(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	return nil
}

type mockErrorStep struct{}

func (s mockErrorStep) Name() string {
	return "mock error step"
}

func (s mockErrorStep) Do(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	return errors.New("forced error")
}

func TestNew(t *testing.T) {
	assert.Equal(t, New(), &Flow{})
}

func TestAdd(t *testing.T) {
	f := &Flow{}
	f.Add(mockStep{}, mockStep{})
	f.Add(mockStep{})
	assert.Equal(t, len(*f), 3)
	for _, step := range *f {
		assert.Equal(t, step.Name(), "mock step")
		assert.Nil(t, step.Do(context.Background(), nil, &http.Request{}))
	}
}

func TestPlay(t *testing.T) {
	f := &Flow{}
	f.Add(mockStep{})
	assert.Nil(t, f.Play(context.Background(), nil, &http.Request{}))
	f.Add(mockErrorStep{})
	sp := mocktracer.New().StartSpan("mock").(*mocktracer.MockSpan)
	assert.Error(t, f.Play(opentracing.ContextWithSpan(context.Background(), sp), nil, &http.Request{}))
	assert.Equal(t, len(sp.Tags()), 1)
	assert.Equal(t, len(sp.Logs()), 1)
}
