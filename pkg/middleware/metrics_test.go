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
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"

	"github.com/northwesternmutual/kanali/pkg/log"
	"github.com/northwesternmutual/kanali/pkg/tags"
)

func TestMetrics(t *testing.T) {
	lvl := zap.NewAtomicLevelAt(zapcore.DebugLevel)
	core, logs := observer.New(lvl)
	defer log.SetLogger(zap.New(core)).Restore()

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/foo", nil)
	msg := "message"
	http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		Metrics(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.WriteHeader(http.StatusTeapot)
			rw.Write([]byte(msg))
		})).ServeHTTP(rw, req)
		assert.Equal(t, 1, logs.FilterMessage("request details").Len())
		assert.Equal(t, 1, logs.FilterField(zap.String(tags.HTTPRequestMethod, "GET")).Len())
		assert.Equal(t, 1, logs.FilterField(zap.Int(tags.HTTPResponseStatusCode, http.StatusTeapot)).Len())
		assert.Equal(t, msg, rec.Body.String())
		assert.Equal(t, http.StatusTeapot, rec.Code)
	}).ServeHTTP(rec, req)
}

func TestGetRequDuration(t *testing.T) {
	t0, t1 := time.Date(2018, time.January, 1, 2, 3, 4, 1000000, time.Local), time.Date(2018, time.January, 1, 2, 3, 4, 256000000, time.Local)
	assert.Equal(t, float64(255), getReqDuration(t0, t1))
}
