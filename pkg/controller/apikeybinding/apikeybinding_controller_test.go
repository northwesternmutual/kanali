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
