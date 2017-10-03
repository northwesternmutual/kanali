package logging

import (
	"context"
	"testing"

	"github.com/northwesternmutual/kanali/config"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestWithContext(t *testing.T) {
	logger := WithContext(nil)
	assert.Nil(t, logger)

	logger = WithContext(context.Background())
	assert.Nil(t, logger)

	viper.SetDefault(config.FlagProcessLogLevel.GetLong(), "foo")
	defer viper.Reset()
	Init(nil)
	logger = WithContext(nil)
	assert.True(t, logger.Core().Enabled(zapcore.InfoLevel))
	assert.False(t, logger.Core().Enabled(zapcore.DebugLevel))

	viper.SetDefault(config.FlagProcessLogLevel.GetLong(), "Warn")
	Init(nil)
	logger = WithContext(context.Background())
	assert.True(t, logger.Core().Enabled(zapcore.WarnLevel))
	assert.False(t, logger.Core().Enabled(zapcore.InfoLevel))

	core, obsvr := observer.New(zap.NewAtomicLevelAt(zapcore.InfoLevel))
	Init(core)
	assert.False(t, logger.Core().Enabled(zapcore.InfoLevel))
	assert.False(t, logger.Core().Enabled(zapcore.DebugLevel))

	ctx := NewContext(context.Background(), zap.String("foo", "bar"))
	logger = WithContext(ctx)
	logger.Info("test log")
	assert.Equal(t, 1, len(obsvr.All()[obsvr.Len()-1].Context))
	assert.Equal(t, "foo", obsvr.All()[obsvr.Len()-1].Context[0].Key)
	assert.Equal(t, "bar", obsvr.All()[obsvr.Len()-1].Context[0].String)
}
