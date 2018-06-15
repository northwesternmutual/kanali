// Copyright (c) 2018 Northwestern Mutual.
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

package middleware

import (
	"net/http"

	newrelic "github.com/newrelic/go-agent"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/northwesternmutual/kanali/pkg/log"
	"github.com/northwesternmutual/kanali/pkg/tags"
)

// newRelicLoggingShim implemnts https://godoc.org/github.com/newrelic/go-agent#Logger
type newRelicLoggingShim struct {
	l *zap.Logger
}

// Error implements https://godoc.org/github.com/newrelic/go-agent#Logger
func (shim *newRelicLoggingShim) Error(msg string, context map[string]interface{}) {
	shim.l.Error(msg, zapifyMap(context)...)
}

// Warn implements https://godoc.org/github.com/newrelic/go-agent#Logger
func (shim *newRelicLoggingShim) Warn(msg string, context map[string]interface{}) {
	shim.l.Warn(msg, zapifyMap(context)...)
}

// Info implements https://godoc.org/github.com/newrelic/go-agent#Logger
func (shim *newRelicLoggingShim) Info(msg string, context map[string]interface{}) {
	shim.l.Info(msg, zapifyMap(context)...)
}

// Debug implements https://godoc.org/github.com/newrelic/go-agent#Logger
func (shim *newRelicLoggingShim) Debug(msg string, context map[string]interface{}) {
	shim.l.Debug(msg, zapifyMap(context)...)
}

// DebugEnabled implements https://godoc.org/github.com/newrelic/go-agent#Logger
func (shim *newRelicLoggingShim) DebugEnabled() bool {
	return shim.l.Core().Enabled(zapcore.DebugLevel)
}

func zapifyMap(m map[string]interface{}) []zapcore.Field {
	fields := make([]zapcore.Field, len(m))
	for k, v := range m {
		fields = append(fields, zap.Any(k, v))
	}
	return fields
}

func NewRelic(enabled bool, licenseKey string) func(next http.Handler) http.Handler {
	logger := log.WithContext(nil)

	config := newrelic.NewConfig(tags.AppName, licenseKey)
	config.Enabled = enabled
	config.Logger = &newRelicLoggingShim{
		l: log.WithContext(nil).Named("New Relic"),
	}

	app, err := newrelic.NewApplication(config)
	if err != nil {
		logger.Warn("error creating New Relic application",
			zap.String(tags.Error, err.Error()),
		)
		return NoOp
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger := log.WithContext(r.Context())
			tnx := app.StartTransaction(r.URL.EscapedPath(), w, r)

			defer func() {
				if err := tnx.End(); err != nil {
					logger.Error("could not end New Relic transaction",
						zap.String(tags.Error, err.Error()),
					)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
