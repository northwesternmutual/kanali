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
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"time"

	"go.uber.org/zap"

	"github.com/northwesternmutual/kanali/pkg/log"
	"github.com/northwesternmutual/kanali/pkg/metrics"
	"github.com/northwesternmutual/kanali/pkg/tags"
	"github.com/northwesternmutual/kanali/pkg/utils"
)

// Metrics is a middleware that will log and report metrics
// corresponding to the current request.
func Metrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		rec := httptest.NewRecorder()

		// This is a workaround until a proper
		// middleware where the errors fall through
		// is implemented.
		for k, v := range w.Header() {
			rec.Header()[k] = v
		}

		metrics.RequestInFlightCount.Inc()

		next.ServeHTTP(rec, r)

		metrics.RequestInFlightCount.Dec()

		endTime := time.Now()

		logger := log.WithContext(r.Context())

		logger.Info("request details",
			zap.String(tags.HTTPRequestRemoteAddress, r.RemoteAddr),
			zap.String(tags.HTTPRequestMethod, r.Method),
			zap.String(tags.HTTPRequestURLPath, utils.ComputeURLPath(r.URL)),
			zap.String(tags.HTTPRequestDuration, fmt.Sprintf("%gms", getReqDuration(startTime, endTime))),
			zap.Int(tags.HTTPResponseStatusCode, rec.Code),
		)

		metrics.RequestLatency.WithLabelValues(strconv.Itoa(rec.Code), r.Method).Observe(getReqDuration(startTime, endTime))

		if err := utils.TransferResponse(rec, w); err != nil {
			logger.Error(fmt.Sprintf("error writing response: %s", err))
		}
	})
}

func getReqDuration(start, finish time.Time) float64 {
	return (float64(finish.UnixNano()-start.UnixNano()) / float64(time.Millisecond))
}
