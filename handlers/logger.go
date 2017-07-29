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

package handlers

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/northwesternmutual/kanali/monitor"
)

// Logger creates a custom http.Handler that logs details around a request
// along with creating contextual request metrics. When the request is complete
// these metrics will be writtin to Influxdb
func Logger(influxCtlr *monitor.InfluxController, inner Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		inner.Context = context.WithValue(context.Background(), monitor.MetricsKey, monitor.New())

		t0 := time.Now()
		inner.serveHTTP(w, r)
		t1 := time.Now()

		logrus.WithFields(logrus.Fields{
			"client ip": strings.Split(r.RemoteAddr, ":")[0],
			"method":    r.Method,
			"uri":       r.RequestURI,
			"totalTime": int(t1.Sub(t0) / time.Millisecond),
		}).Info("request details")

		inner.Context = monitor.AddCtxMetric(inner.Context, "total_time", strconv.Itoa(int(t1.Sub(t0)/time.Millisecond)))
		inner.Context = monitor.AddCtxMetric(inner.Context, "http_method", r.Method)
		inner.Context = monitor.AddCtxMetric(inner.Context, "http_uri", r.RequestURI)
		inner.Context = monitor.AddCtxMetric(inner.Context, "client_ip", strings.Split(r.RemoteAddr, ":")[0])

		if err := influxCtlr.WriteRequestData(inner.Context); err != nil {
			logrus.Warnf(err.Error())
		} else {
			logrus.Infof("successfully wrote request details to influxdb")
		}

	})

}
