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
