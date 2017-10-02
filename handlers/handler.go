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
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/northwesternmutual/kanali/logging"
	"github.com/northwesternmutual/kanali/metrics"
	"github.com/northwesternmutual/kanali/monitor"
	"github.com/northwesternmutual/kanali/spec"
	"github.com/northwesternmutual/kanali/tracer"
	"github.com/northwesternmutual/kanali/utils"
	"github.com/opentracing/opentracing-go"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
)

// Handler is used to provide additional parameters to an HTTP handler
type Handler struct {
	*monitor.InfluxController
	H func(ctx context.Context, proxy *spec.APIProxy, m *metrics.Metrics, w http.ResponseWriter, r *http.Request, trace opentracing.Span) error
}

// ServeHTTP serves an HTTP request
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	httpContext := context.Background()
	rqCtx := logging.NewContext(httpContext, zap.Stringer("correlation_id", uuid.NewV4()))
	logger := logging.WithContext(rqCtx)

	t0 := time.Now()
	m := &metrics.Metrics{}

	normalize(r)

	defer func() {
		m.Add(
			metrics.Metric{Name: "total_time", Value: int(time.Now().Sub(t0) / time.Millisecond), Index: false},
			metrics.Metric{Name: "http_method", Value: r.Method, Index: true},
			metrics.Metric{Name: "http_uri", Value: r.URL.EscapedPath(), Index: false},
			metrics.Metric{Name: "client_ip", Value: strings.Split(r.RemoteAddr, ":")[0], Index: false},
		)
		logger.Info("request details",
			zap.String(tracer.HTTPRequestRemoteAddress, strings.Split(r.RemoteAddr, ":")[0]),
			zap.String(tracer.HTTPRequestMethod, r.Method),
			zap.String(tracer.HTTPRequestURLPath, r.URL.EscapedPath()),
		)
		go func() {
			if err := h.InfluxController.WriteRequestData(m); err != nil {
				logger.Warn(err.Error())
			} else {
				logger.Debug("wrote metrics to InfluxDB")
			}
		}()
	}()

	sp := opentracing.StartSpan(fmt.Sprintf("%s %s",
		r.Method,
		r.URL.EscapedPath(),
	))
	defer sp.Finish()

	tracer.HydrateSpanFromRequest(r, sp)

	err := h.H(httpContext, &spec.APIProxy{}, m, w, r, sp)
	if err == nil {
		return
	}

	// all errors will need the application/json Content-Type header
	w.Header().Set("Content-Type", "application/json")

	// we'll have multiple types off errors
	switch e := err.(type) {
	case utils.Error:

		sp.SetTag(tracer.HTTPResponseStatusCode, e.Status())

		logger.Info(err.Error(),
			zap.String(tracer.HTTPRequestMethod, r.Method),
			zap.String(tracer.HTTPRequestURLPath, r.URL.EscapedPath()),
		)

		m.Add(metrics.Metric{Name: "http_response_code", Value: strconv.Itoa(e.Status()), Index: true})

		errStatus, err := json.Marshal(utils.JSONErr{Code: e.Status(), Msg: e.Error()})
		if err != nil {
			logger.Warn(err.Error())
		} else {
			sp.SetTag(tracer.HTTPResponseBody, string(errStatus))
		}

		// write error code to response
		w.WriteHeader(e.Status())

		// write error message to response
		if err := json.NewEncoder(w).Encode(utils.JSONErr{Code: e.Status(), Msg: e.Error()}); err != nil {
			logger.Error(err.Error())
		}

	default:

		sp.SetTag(tracer.HTTPResponseStatusCode, http.StatusInternalServerError)

		logger.Info("unknown error",
			zap.String(tracer.HTTPRequestMethod, r.Method),
			zap.String(tracer.HTTPRequestURLPath, r.URL.EscapedPath()),
		)

		m.Add(metrics.Metric{Name: "http_response_code", Value: strconv.Itoa(http.StatusInternalServerError), Index: true})

		errStatus, err := json.Marshal(utils.JSONErr{Code: http.StatusInternalServerError, Msg: "unknown error"})
		if err != nil {
			logger.Warn(err.Error())
		} else {
			sp.SetTag(tracer.HTTPResponseBody, string(errStatus))
		}

		// write error code to response
		w.WriteHeader(http.StatusInternalServerError)

		// write error message to response
		if err := json.NewEncoder(w).Encode(utils.JSONErr{Code: http.StatusInternalServerError, Msg: "unknown error"}); err != nil {
			logger.Error(err.Error())
		}

	}
}

func normalize(r *http.Request) {
	r.URL.Path = utils.NormalizePath(r.URL.Path)
	r.URL.RawPath = r.URL.EscapedPath()
}
