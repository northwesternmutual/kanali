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

package apikeybinding

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"

	"github.com/northwesternmutual/kanali/pkg/log"
	store "github.com/northwesternmutual/kanali/pkg/store/kanali/v2"
	"github.com/northwesternmutual/kanali/pkg/tags"
	"github.com/northwesternmutual/kanali/test/builder"
)

func TestOnAdd(t *testing.T) {
	core, logs := observer.New(zap.NewAtomicLevelAt(zapcore.DebugLevel))
	defer log.SetLogger(zap.New(core)).Restore()
	defer store.ApiKeyBindingStore().Clear()

	apikeybinding := builder.NewApiKeyBinding("foo", "bar").WithKeys(
		builder.NewKeyAccess("frank").WithDefaultRule(
			builder.NewRule().WithGlobal().NewOrDie(),
		).NewOrDie(),
	).NewOrDie()

	NewController().OnAdd(apikeybinding)
	assert.NotNil(t, store.ApiKeyBindingStore().GetHightestPriorityRule("bar", "foo", "frank", "/foo"))
	assert.Equal(t, 1, logs.FilterField(zap.String(tags.KanaliApiKeyBindingName, "foo")).FilterField(zap.String(tags.KanaliApiKeyBindingNamespace, "bar")).Len())

	NewController().OnAdd(nil)
	assert.Equal(t, 1, logs.FilterMessageSnippet("malformed").Len())
}

func TestOnUpdate(t *testing.T) {
	core, logs := observer.New(zap.NewAtomicLevelAt(zapcore.DebugLevel))
	defer log.SetLogger(zap.New(core)).Restore()
	defer store.ApiKeyBindingStore().Clear()

	apikeybindingOld := builder.NewApiKeyBinding("foo", "bar").WithKeys(
		builder.NewKeyAccess("frank").WithDefaultRule(
			builder.NewRule().WithGlobal().NewOrDie(),
		).NewOrDie(),
	).NewOrDie()
	apikeybindingNew := builder.NewApiKeyBinding("foo", "bar").WithKeys(
		builder.NewKeyAccess("joe").WithDefaultRule(
			builder.NewRule().WithGlobal().NewOrDie(),
		).NewOrDie(),
	).NewOrDie()

	NewController().OnUpdate(nil, nil)
	assert.Equal(t, 1, logs.FilterMessageSnippet("malformed").Len())
	NewController().OnUpdate(nil, "foo")
	assert.Equal(t, 2, logs.FilterMessageSnippet("malformed").Len())
	NewController().OnUpdate(apikeybindingOld, "foo")
	assert.Equal(t, 3, logs.FilterMessageSnippet("malformed").Len())

	store.ApiKeyBindingStore().Set(apikeybindingOld)
	NewController().OnUpdate(apikeybindingOld, apikeybindingNew)
	assert.Nil(t, store.ApiKeyBindingStore().GetHightestPriorityRule("bar", "foo", "frank", "/foo"))
	assert.NotNil(t, store.ApiKeyBindingStore().GetHightestPriorityRule("bar", "foo", "joe", "/foo"))
}

func TestOnDelete(t *testing.T) {
	core, logs := observer.New(zap.NewAtomicLevelAt(zapcore.DebugLevel))
	defer log.SetLogger(zap.New(core)).Restore()
	defer store.ApiKeyBindingStore().Clear()

	apikeybinding := builder.NewApiKeyBinding("foo", "bar").WithKeys(
		builder.NewKeyAccess("frank").WithDefaultRule(
			builder.NewRule().WithGlobal().NewOrDie(),
		).NewOrDie(),
	).NewOrDie()

	NewController().OnDelete(nil)
	assert.Equal(t, 1, logs.FilterMessageSnippet("malformed").Len())

	store.ApiKeyBindingStore().Set(apikeybinding)
	NewController().OnDelete(apikeybinding)
	assert.Nil(t, store.ApiKeyBindingStore().GetHightestPriorityRule("bar", "foo", "frank", "/foo"))
	assert.Equal(t, 1, logs.FilterMessageSnippet("deleted").Len())
	assert.Equal(t, 1, logs.FilterField(zap.String(tags.KanaliApiKeyBindingName, "foo")).FilterField(zap.String(tags.KanaliApiKeyBindingNamespace, "bar")).Len())
}
