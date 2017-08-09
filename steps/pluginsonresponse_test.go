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
	"testing"

	"github.com/northwesternmutual/kanali/controller"
	"github.com/northwesternmutual/kanali/spec"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/stretchr/testify/assert"
)

func TestPluginsOnResponseGetName(t *testing.T) {
	step := PluginsOnResponseStep{}
	assert.Equal(t, step.GetName(), "Plugin OnResponse", "step name is incorrect")
}

func TestDoOnResponse(t *testing.T) {
	assert.Equal(t, doOnResponse(context.Background(), spec.APIProxy{}, controller.Controller{}, nil, nil, opentracing.StartSpan("test span"), fakePanicPlugin{}).Error(), "OnResponse paniced")
  assert.Equal(t, doOnResponse(context.Background(), spec.APIProxy{}, controller.Controller{}, nil, nil, opentracing.StartSpan("test span"), fakeErrorPlugin{}).Error(), "error")
  assert.Nil(t, doOnResponse(context.Background(), spec.APIProxy{}, controller.Controller{}, nil, nil, opentracing.StartSpan("test span"), fakeSuccessPlugin{}))
}
