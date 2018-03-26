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
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestInit(t *testing.T) {
	assert.Equal(t, zapcore.InfoLevel, wrappedLogger.level.Level())
	assert.True(t, wrappedLogger.zap.Core().Enabled(zapcore.InfoLevel))
	assert.False(t, wrappedLogger.zap.Core().Enabled(zapcore.DebugLevel))
}

func TestSetLevel(t *testing.T) {
	assert.False(t, wrappedLogger.zap.Core().Enabled(zapcore.DebugLevel))
	lvl := DebugLevel
	SetLevel(&lvl)
	assert.True(t, wrappedLogger.zap.Core().Enabled(zapcore.DebugLevel))
}

func TestNewContext(t *testing.T) {
	lvl := zap.NewAtomicLevelAt(zapcore.DebugLevel)
	core, logs := observer.New(lvl)
	defer SetLogger(zap.New(core)).Restore()

	ctx := NewContext(context.Background(), zap.String("foo", "bar"))
	logger, ok := ctx.Value(loggerKey).(*zap.Logger)
	assert.True(t, ok)
	logger.Debug("foo")
	assert.Equal(t, 1, logs.FilterField(zap.String("foo", "bar")).Len())
}

func TestWithContext(t *testing.T) {
	assert.Equal(t, wrappedLogger.zap, WithContext(nil))
	assert.Equal(t, wrappedLogger.zap, WithContext(context.Background()))
	ctx := NewContext(context.Background(), zap.String("foo", "bar"))
	logger, ok := ctx.Value(loggerKey).(*zap.Logger)
	assert.True(t, ok)
	assert.Equal(t, logger, WithContext(ctx))
}

// There is no way actually test this yet.
// But I still want the code coverage :)
// https://github.com/prometheus/client_golang/issues/58
func TestAddMetrics(t *testing.T) {
	addMetrics(zapcore.Entry{Level: zap.DebugLevel})
	addMetrics(zapcore.Entry{Level: zap.ErrorLevel})
	addMetrics(zapcore.Entry{Level: zap.FatalLevel})
	addMetrics(zapcore.Entry{Level: zap.InfoLevel})
	addMetrics(zapcore.Entry{Level: zap.WarnLevel})
	addMetrics(zapcore.Entry{Level: zap.PanicLevel})
	addMetrics(zapcore.Entry{Level: zap.DPanicLevel})
}
