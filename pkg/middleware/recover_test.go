package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"

	"github.com/northwesternmutual/kanali/pkg/log"
)

func TestRecover(t *testing.T) {
	lvl := zap.NewAtomicLevelAt(zapcore.DebugLevel)
	core, logs := observer.New(lvl)
	defer log.SetLogger(zap.New(core)).Restore()

	req, _ := http.NewRequest("GET", "/foo", nil)
	rec := httptest.NewRecorder()

	assert.NotPanics(t, assert.PanicTestFunc(func() {
		Recover(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			panic("help me")
		})).ServeHTTP(rec, req)
	}))

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Equal(t, 1, logs.Len())
}
