package logging

import (
	"context"
	"os"

	"github.com/northwesternmutual/kanali/config"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type loggerKeyType int

const (
	loggerKey loggerKeyType = iota
)

var logger *zap.Logger

// Init instantiates a global logger that Kanali will use
func Init(core zapcore.Core) {
	if core != nil {
		logger = zap.New(core)
		return
	}
	// create new logger with default production configuration
	cfg := zap.NewProductionConfig()
	// create new logging level
	level := zap.NewAtomicLevel()
	// attempt to set the logging level
	level.UnmarshalText([]byte(viper.GetString(config.FlagProcessLogLevel.GetLong())))
	// create new logger from above configuration
	l, err := cfg.Build()
	if err == nil {
		cfg.Level.SetLevel(level.Level())
		logger = l
		return
	}
	// if error creating logger from above configuration, create a backup logger
	logger = zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(cfg.EncoderConfig),
		zapcore.Lock(os.Stderr),
		level,
	))
}

// NewContext creates a new context the given contextual fields
func NewContext(ctx context.Context, fields ...zapcore.Field) context.Context {
	return context.WithValue(ctx, loggerKey, WithContext(ctx).With(fields...))
}

// WithContext returns a logger from the given context
func WithContext(ctx context.Context) *zap.Logger {
	if ctx == nil {
		return logger
	}
	if ctxLogger, ok := ctx.Value(loggerKey).(*zap.Logger); ok {
		return ctxLogger
	}
	return logger
}
