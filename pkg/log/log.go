// Copyright (c) 2017 Northwestern Mutual.
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

package log

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/northwesternmutual/kanali/pkg/metrics"
)

const loggerKey = iota

type logger struct {
	zap   *zap.Logger
	level zap.AtomicLevel
}

var wrappedLogger *logger

func init() {
	wrappedLogger = new(logger)
	wrappedLogger.level = zap.NewAtomicLevel()

	config := zap.NewProductionConfig()
	config.Level = wrappedLogger.level
	l, err := config.Build(
		zap.Hooks(addMetrics),
		zap.AddCaller(),
		zap.AddStacktrace(zap.ErrorLevel),
	)
	if err != nil {
		panic(err)
	}

	wrappedLogger.zap = l
}

func SetLevel(lvl string) {
	wrappedLogger.level.UnmarshalText([]byte(lvl))
}

// NewContext creates a new context the given contextual fields
func NewContext(ctx context.Context, fields ...zapcore.Field) context.Context {
	return context.WithValue(ctx, loggerKey, WithContext(ctx).With(fields...))
}

// WithContext returns a logger from the given context
func WithContext(ctx context.Context) *zap.Logger {
	if ctx == nil {
		return wrappedLogger.zap
	}
	if ctxLogger, ok := ctx.Value(loggerKey).(*zap.Logger); ok {
		return ctxLogger
	}
	return wrappedLogger.zap
}

func addMetrics(e zapcore.Entry) error {
	switch e.Level {
	case zap.DebugLevel:
		metrics.LoggingCount.WithLabelValues("debug").Inc()
	case zap.ErrorLevel:
		metrics.LoggingCount.WithLabelValues("error").Inc()
	case zap.FatalLevel:
		metrics.LoggingCount.WithLabelValues("fatal").Inc()
	case zap.InfoLevel:
		metrics.LoggingCount.WithLabelValues("info").Inc()
	case zap.WarnLevel:
		metrics.LoggingCount.WithLabelValues("warn").Inc()
	case zap.PanicLevel:
		metrics.LoggingCount.WithLabelValues("panic").Inc()
	}
	return nil
}
