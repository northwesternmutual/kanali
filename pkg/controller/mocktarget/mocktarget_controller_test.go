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
		*builder.NewRoute("/foo").WithMethods("GET").NewOrDie(),
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
		*builder.NewRoute("/foo").WithMethods("GET").NewOrDie(),
	).NewOrDie()
	mocktargetNew := builder.NewMockTarget("foo", "bar").WithRoute(
		*builder.NewRoute("/bar").WithMethods("GET").NewOrDie(),
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
		*builder.NewRoute("/foo").WithMethods("GET").NewOrDie(),
	).NewOrDie()

	NewController().OnDelete(nil)
	assert.Equal(t, 1, logs.FilterMessageSnippet("malformed").Len())

	store.MockTargetStore().Set(mocktarget)
	assert.NotNil(t, store.MockTargetStore().Get("bar", "foo", "/foo", "GET"))
	NewController().OnDelete(mocktarget)
	assert.Nil(t, store.MockTargetStore().Get("bar", "foo", "/foo", "GET"))
	assert.Equal(t, 1, logs.FilterMessageSnippet("deleted").Len())

}
