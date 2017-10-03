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

package tracer

import (
	"testing"

	"github.com/northwesternmutual/kanali/logging"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestCustomLoggerError(t *testing.T) {
	core, obsvr := observer.New(zap.NewAtomicLevelAt(zapcore.DebugLevel))
	logging.Init(core)
	logger := customLogger{}

	logger.Error("custom error message")
	assert.Equal(t, zapcore.ErrorLevel, obsvr.All()[obsvr.Len()-1].Entry.Level)
	assert.Equal(t, "custom error message", obsvr.All()[obsvr.Len()-1].Entry.Message)
	assert.Equal(t, 0, len(obsvr.All()[obsvr.Len()-1].Context))

	logger.Infof("custom %s message", "info")
	assert.Equal(t, zapcore.InfoLevel, obsvr.All()[obsvr.Len()-1].Entry.Level)
	assert.Equal(t, "custom info message", obsvr.All()[obsvr.Len()-1].Entry.Message)
	assert.Equal(t, 0, len(obsvr.All()[obsvr.Len()-1].Context))
}
