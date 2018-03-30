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
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/northwesternmutual/kanali/pkg/metrics"
)

const loggerKey = iota

type logger struct {
	zap   *zap.Logger
	level zap.AtomicLevel
}

type tempLogger struct {
	previous zap.Logger
}

var wrappedLogger *logger

func init() {
	// Create a new logging configuration. Using zap's default production config.
	// More details here: https://github.com/uber-go/zap/blob/master/config.go#L114-L127
	// If further tweaking is needed, either this one can be tweaked ore one
	// can be created from scratch. Note that the default level is info.
	config := zap.NewProductionConfig()
	l, err := config.Build(
		zap.Hooks(addMetrics),
		zap.AddCaller(),
		zap.AddStacktrace(zap.ErrorLevel),
	)
	if err != nil {
		panic(err)
	}

	wrappedLogger = &logger{
		zap:   l.Named("kanali"),
		level: config.Level,
	}
}

// SetLogger mutates the global logger. This function is often
// used in conjunction with Restore() which restores the
// global logger to the one used before a call to this function.
func SetLogger(l *zap.Logger) *tempLogger {
	tmp := &tempLogger{
		previous: *(wrappedLogger.zap),
	}
	wrappedLogger.zap = l
	return tmp
}

// Restore mutates the global logger and sets it to the
// previous logger.
func (t *tempLogger) Restore() {
	wrappedLogger.zap = &(t.previous)
}

// SetLevel dynamically sets the logging level.
func SetLevel(lvl *Level) {
	if err := wrappedLogger.level.UnmarshalText([]byte(lvl.String())); err != nil {
		wrappedLogger.zap.Warn(fmt.Sprintf("error setting lot level: %s", err))
	}
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
	if e.Level == zap.DPanicLevel {
		e.Level = zap.PanicLevel
	}
	metrics.LoggingCount.WithLabelValues(e.Level.String()).Inc()
	return nil
}
