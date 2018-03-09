package server

import (
	"os"
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"

	"github.com/northwesternmutual/kanali/pkg/log"
)

func TestSetupSignalHandler(t *testing.T) {
	SetupSignalHandler()
	assert.Panics(t, assert.PanicTestFunc(func() {
		SetupSignalHandler()
	}))
}

func TestLogSignal(t *testing.T) {
	lvl := zap.NewAtomicLevelAt(zapcore.DebugLevel)
	core, logs := observer.New(lvl)
	defer log.SetLogger(zap.New(core)).Restore()

	logSignal(os.Interrupt)
	logSignal(syscall.SIGTERM)
	assert.Equal(t, 1, logs.FilterMessageSnippet("SIGTERM").Len())
	assert.Equal(t, 1, logs.FilterMessageSnippet("SIGINT").Len())
}
