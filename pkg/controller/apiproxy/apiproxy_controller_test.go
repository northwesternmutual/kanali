package apiproxy

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
	defer store.ApiProxyStore().Clear()

	NewController().OnAdd(builder.NewApiProxy("foo", "bar").WithSourcePath("/foo").WithTargetBackendEndpoint("http://foo.bar.com").NewOrDie())
	assert.NotNil(t, store.ApiProxyStore().Get("/foo"))
	assert.Equal(t, 1, logs.FilterField(zap.String(tags.KanaliProxyName, "foo")).FilterField(zap.String(tags.KanaliProxyNamespace, "bar")).Len())

	NewController().OnAdd(nil)
	assert.Equal(t, 1, logs.FilterMessageSnippet("malformed").Len())
}

func TestOnUpdate(t *testing.T) {
	core, logs := observer.New(zap.NewAtomicLevelAt(zapcore.DebugLevel))
	defer log.SetLogger(zap.New(core)).Restore()
	defer store.ApiProxyStore().Clear()

	apiproxyOld := builder.NewApiProxy("foo", "bar").WithSourcePath("/foo").WithTargetBackendEndpoint("http://foo.bar.com").NewOrDie()
	apiproxyNew := builder.NewApiProxy("foo", "bar").WithSourcePath("/foo").WithTargetPath("/bar").WithTargetBackendEndpoint("http://foo.bar.com").NewOrDie()

	NewController().OnUpdate(nil, nil)
	assert.Equal(t, 1, logs.FilterMessageSnippet("malformed").Len())
	NewController().OnUpdate(nil, "foo")
	assert.Equal(t, 2, logs.FilterMessageSnippet("malformed").Len())
	NewController().OnUpdate(apiproxyOld, "foo")
	assert.Equal(t, 3, logs.FilterMessageSnippet("malformed").Len())

	store.ApiProxyStore().Set(apiproxyOld)
	assert.Equal(t, "/", store.ApiProxyStore().Get("/foo").Spec.Target.Path)
	NewController().OnUpdate(apiproxyOld, apiproxyNew)
	assert.Equal(t, "/bar", store.ApiProxyStore().Get("/foo").Spec.Target.Path)
	assert.Equal(t, 1, logs.FilterMessageSnippet("updated").Len())
	assert.Equal(t, 1, logs.FilterField(zap.String(tags.KanaliProxyName, "foo")).FilterField(zap.String(tags.KanaliProxyNamespace, "bar")).Len())
}

func TestOnDelete(t *testing.T) {
	core, logs := observer.New(zap.NewAtomicLevelAt(zapcore.DebugLevel))
	defer log.SetLogger(zap.New(core)).Restore()
	defer store.ApiProxyStore().Clear()

	apiproxy := builder.NewApiProxy("foo", "bar").WithSourcePath("/foo").WithTargetBackendEndpoint("http://foo.bar.com").NewOrDie()
	store.ApiProxyStore().Set(apiproxy)

	NewController().OnDelete(nil)
	assert.Equal(t, 1, logs.FilterMessageSnippet("malformed").Len())

	assert.NotNil(t, store.ApiProxyStore().Get("/foo"))
	NewController().OnDelete(apiproxy)
	assert.Nil(t, store.ApiProxyStore().Get("/foo"))
	assert.Equal(t, 1, logs.FilterMessageSnippet("deleted").Len())
	assert.Equal(t, 1, logs.FilterField(zap.String(tags.KanaliProxyName, "foo")).FilterField(zap.String(tags.KanaliProxyNamespace, "bar")).Len())
}
