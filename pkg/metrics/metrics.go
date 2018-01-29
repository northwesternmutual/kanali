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

package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	RequestLatency = prometheus.NewHistogramVec(prometheus.HistogramOpts{
    Namespace: "http",
		Name: "request_latency_milliseconds",
		Help: "Latency of HTTP requests.",
	}, []string{"method"})
  RequestCount = prometheus.NewCounterVec(prometheus.CounterOpts{
    Namespace: "http",
		Name: "request_count_total",
		Help: "Count of all HTTP requests.",
	}, []string{"code", "method"})
  RequestErrorCount = prometheus.NewCounterVec(prometheus.CounterOpts{
    Namespace: "http",
		Name: "request_error_count_total",
		Help: "Count of all HTTP requests resulting in an error.",
	}, []string{"code", "method"})
  RequestInFlightCount = prometheus.NewGauge(prometheus.GaugeOpts{
    Namespace: "http",
		Name: "request_in_flight_total",
		Help: "Count of all HTTP requests still in flight.",
	})

  LoggingInfoTotal = prometheus.NewCounter(prometheus.CounterOpts{
    Namespace: "logging",
		Name: "info_total",
		Help: "Number of info level logs.",
	})
  LoggingDebugTotal = prometheus.NewCounter(prometheus.CounterOpts{
    Namespace: "logging",
		Name: "debug_total",
		Help: "Number of debug level logs.",
	})
  LoggingWarnTotal = prometheus.NewCounter(prometheus.CounterOpts{
    Namespace: "logging",
		Name: "warn_total",
		Help: "Number of warn level logs.",
	})
  LoggingErrorTotal = prometheus.NewCounter(prometheus.CounterOpts{
    Namespace: "logging",
		Name: "error_total",
		Help: "Number of error level logs.",
	})
  LoggingPanicTotal = prometheus.NewCounter(prometheus.CounterOpts{
    Namespace: "logging",
		Name: "panic_total",
		Help: "Number of panic level logs.",
	})
  LoggingFatalTotal = prometheus.NewCounter(prometheus.CounterOpts{
    Namespace: "logging",
		Name: "fatal_total",
		Help: "Number of fatal level logs.",
	})
  LoggingTotal = prometheus.NewCounter(prometheus.CounterOpts{
    Namespace: "logging",
		Name: "total",
		Help: "Number of level logs.",
	})
)

func init() {
	prometheus.MustRegister(
    RequestLatency,
    RequestCount,
    RequestErrorCount,
    RequestInFlightCount,

    LoggingInfoTotal,
    LoggingDebugTotal,
    LoggingWarnTotal,
    LoggingErrorTotal,
    LoggingPanicTotal,
    LoggingFatalTotal,
    LoggingTotal,
  )
}
