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
	"net/http"
	"net/url"
	"testing"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/mocktracer"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"

	"github.com/northwesternmutual/kanali/pkg/errors"
	"github.com/northwesternmutual/kanali/pkg/log"
	store "github.com/northwesternmutual/kanali/pkg/store/kanali/v2"
	"github.com/northwesternmutual/kanali/pkg/tags"
	"github.com/northwesternmutual/kanali/test/builder"
)

func TestValidateProxyName(t *testing.T) {
	assert.Equal(t, "Validate Proxy", ValidateProxyStep().Name())
}

func TestValidateProxyDo(t *testing.T) {
	lvl := zap.NewAtomicLevelAt(zapcore.DebugLevel)
	core, logs := observer.New(lvl)
	defer log.SetLogger(zap.New(core)).Restore()

	apiproxy := builder.NewApiProxy("foo", "bar").WithSourcePath("/foo").WithTargetBackendEndpoint("http://foo.bar.com").NewOrDie()
	store.ApiProxyStore().Set(apiproxy)
	u, _ := url.Parse("/foo/bar")
	sp := mocktracer.New().StartSpan("mock").(*mocktracer.MockSpan)
	assert.Nil(t, ValidateProxyStep().Do(opentracing.ContextWithSpan(context.Background(), sp), nil, &http.Request{
		URL: u,
	}))
	assert.Equal(t, 1, logs.Len())
	assert.Equal(t, 1, logs.FilterField(zap.String(tags.KanaliProxyName, "foo")).FilterField(zap.String(tags.KanaliProxyNamespace, "bar")).Len())
	assert.Equal(t, "foo", sp.Tag(tags.KanaliProxyName))
	assert.Equal(t, "bar", sp.Tag(tags.KanaliProxyNamespace))

	u, _ = url.Parse("/bar/foo")
	err := ValidateProxyStep().Do(context.Background(), nil, &http.Request{
		URL: u,
	})
	assert.Error(t, err)
	assert.Equal(t, errors.ErrorProxyNotFound, err)
	assert.Equal(t, 2, logs.Len())
}
