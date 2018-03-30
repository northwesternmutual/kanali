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

package mocktarget

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"

	"github.com/northwesternmutual/kanali/pkg/log"
	store "github.com/northwesternmutual/kanali/pkg/store/kanali/v2"
	"github.com/northwesternmutual/kanali/test/builder"
)

func TestOnAdd(t *testing.T) {
	core, logs := observer.New(zap.NewAtomicLevelAt(zapcore.DebugLevel))
	defer log.SetLogger(zap.New(core)).Restore()
	defer store.MockTargetStore().Clear()

	mocktarget := builder.NewMockTarget("foo", "bar").WithRoute(
		builder.NewRoute("/foo").WithMethods("GET").NewOrDie(),
	).NewOrDie()

	NewController().OnAdd(nil)
	assert.Equal(t, 1, logs.FilterMessageSnippet("malformed").Len())

	NewController().OnAdd(mocktarget)
	assert.NotNil(t, store.MockTargetStore().Get("bar", "foo", "/foo", "GET"))
	assert.Equal(t, 1, logs.FilterMessageSnippet("added").Len())
}

func TestOnUpdate(t *testing.T) {
	core, logs := observer.New(zap.NewAtomicLevelAt(zapcore.DebugLevel))
	defer log.SetLogger(zap.New(core)).Restore()
	defer store.MockTargetStore().Clear()

	mocktargetOld := builder.NewMockTarget("foo", "bar").WithRoute(
		builder.NewRoute("/foo").WithMethods("GET").NewOrDie(),
	).NewOrDie()
	mocktargetNew := builder.NewMockTarget("foo", "bar").WithRoute(
		builder.NewRoute("/bar").WithMethods("GET").NewOrDie(),
	).NewOrDie()

	NewController().OnUpdate(nil, nil)
	assert.Equal(t, 1, logs.FilterMessageSnippet("malformed").Len())
	NewController().OnUpdate(nil, "foo")
	assert.Equal(t, 2, logs.FilterMessageSnippet("malformed").Len())
	NewController().OnUpdate(mocktargetOld, "foo")
	assert.Equal(t, 3, logs.FilterMessageSnippet("malformed").Len())

	store.MockTargetStore().Set(mocktargetOld)
	NewController().OnUpdate(mocktargetOld, mocktargetNew)
	assert.NotNil(t, store.MockTargetStore().Get("bar", "foo", "/bar", "GET"))
	assert.Nil(t, store.MockTargetStore().Get("bar", "foo", "/foo", "GET"))
}

func TestOnDelete(t *testing.T) {
	core, logs := observer.New(zap.NewAtomicLevelAt(zapcore.DebugLevel))
	defer log.SetLogger(zap.New(core)).Restore()
	defer store.MockTargetStore().Clear()

	mocktarget := builder.NewMockTarget("foo", "bar").WithRoute(
		builder.NewRoute("/foo").WithMethods("GET").NewOrDie(),
	).NewOrDie()

	NewController().OnDelete(nil)
	assert.Equal(t, 1, logs.FilterMessageSnippet("malformed").Len())

	store.MockTargetStore().Set(mocktarget)
	assert.NotNil(t, store.MockTargetStore().Get("bar", "foo", "/foo", "GET"))
	NewController().OnDelete(mocktarget)
	assert.Nil(t, store.MockTargetStore().Get("bar", "foo", "/foo", "GET"))
	assert.Equal(t, 1, logs.FilterMessageSnippet("deleted").Len())

}
