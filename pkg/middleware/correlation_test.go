package middleware

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"

	"github.com/northwesternmutual/kanali/pkg/log"
	"github.com/northwesternmutual/kanali/pkg/tags"
)

func TestCorrelation(t *testing.T) {
	lvl := zap.NewAtomicLevelAt(zapcore.DebugLevel)
	core, logs := observer.New(lvl)
	defer log.SetLogger(zap.New(core)).Restore()

	req, _ := http.NewRequest("GET", "/foo", nil)
	Correlation(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, 1, len(logs.All()))
		assert.Equal(t, tags.HTTPRequestCorrelationId, logs.All()[0].Context[0].Key)
	})).ServeHTTP(nil, req)
}
