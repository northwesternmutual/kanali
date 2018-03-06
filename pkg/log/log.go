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

type tempLogger struct {
	previous zap.Logger
}

type Level zapcore.Level

const (
	// DebugLevel logs are typically voluminous, and are usually disabled in
	// production.
	DebugLevel Level = Level(zapcore.DebugLevel)
	// InfoLevel is the default logging priority.
	InfoLevel Level = Level(zapcore.InfoLevel)
	// WarnLevel logs are more important than Info, but don't need individual
	// human review.
	WarnLevel Level = Level(zapcore.WarnLevel)
	// ErrorLevel logs are high-priority. If an application is running smoothly,
	// it shouldn't generate any error-level logs.
	ErrorLevel Level = Level(zapcore.ErrorLevel)
	// DPanicLevel logs are particularly important errors. In development the
	// logger panics after writing the message.
	DPanicLevel Level = Level(zapcore.DPanicLevel)
	// PanicLevel logs a message, then panics.
	PanicLevel Level = Level(zapcore.PanicLevel)
	// FatalLevel logs a message, then calls os.Exit(1).
	FatalLevel Level = Level(zapcore.FatalLevel)
)

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
		zap:   l,
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
func SetLevel(lvl Level) {
	wrappedLogger.level.SetLevel(zapcore.Level(lvl))
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
