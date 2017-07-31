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

	"github.com/Sirupsen/logrus"
	"github.com/northwesternmutual/kanali/config"
	"github.com/northwesternmutual/kanali/controller"
	"github.com/northwesternmutual/kanali/monitor"
	"github.com/northwesternmutual/kanali/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/spf13/viper"
)

// Handler is used to provide additional parameters to an HTTP handler
type Handler struct {
	*controller.Controller
	context.Context
	H func(ctx context.Context, c *controller.Controller, w http.ResponseWriter, r *http.Request, trace opentracing.Span) error
}

func (h Handler) serveHTTP(w http.ResponseWriter, r *http.Request) {

	// start a global trace
	sp := opentracing.StartSpan(fmt.Sprintf("%s %s",
		r.Method,
		r.URL.Path,
	))

	closer, str, err := utils.DupReaderAndString(r.Body)
	if err != nil {
		logrus.Errorf("error copying request body, data may not be as expected: %s", err.Error())
	}

	// copy new reader into freshly drained reader
	r.Body = closer

	sp.SetTag("http.request_body", str)
	sp.SetTag("http.url", r.RequestURI)
	sp.SetTag("http.method", r.Method)

	jsonHeaders, err := json.Marshal(utils.FlattenHTTPHeaders(utils.OmitHeaderValues(r.Header, viper.GetString(config.FlagHeaderMaskValue.GetLong()), config.FlagApikeyHeaderKey.GetLong())))
	if err != nil {
		logrus.Warnf("could not marsah request headers into JSON - tracing data maybe not be as expected")
	} else {
		sp.SetTag("http.headers", string(jsonHeaders))
	}

	defer sp.Finish()

	err = h.H(h.Context, h.Controller, w, r, sp)

	// handle request errors
	if err != nil {

		// all errors will need the application/json Content-Type header
		w.Header().Set("Content-Type", "application/json")

		// we'll have multiple types off errors
		switch e := err.(type) {
		case utils.Error:

			sp.SetTag("http.status_code", e.Status())

			// log error
			logrus.WithFields(logrus.Fields{
				"method": r.Method,
				"uri":    r.RequestURI,
			}).Error(e.Error())

			h.Context = monitor.AddCtxMetric(h.Context, "http_response_code", strconv.Itoa(e.Status()))

			errStatus, err := json.Marshal(utils.JSONErr{Code: e.Status(), Msg: e.Error()})
			if err != nil {
				logrus.Warnf("could not marsah request headers into JSON - tracing data maybe not be as expected")
			} else {
				sp.SetTag("http.response_body", string(errStatus))
			}

			// write error code to response
			w.WriteHeader(e.Status())

			// write error message to response
			if err := json.NewEncoder(w).Encode(utils.JSONErr{Code: e.Status(), Msg: e.Error()}); err != nil {
				logrus.Fatal(err.Error())
			}

		default:

			sp.SetTag("http.status_code", http.StatusInternalServerError)

			// log error
			logrus.WithFields(logrus.Fields{
				"method": r.Method,
				"uri":    r.RequestURI,
			}).Error("unknown error")

			h.Context = monitor.AddCtxMetric(h.Context, "http_response_code", strconv.Itoa(http.StatusInternalServerError))

			errStatus, err := json.Marshal(utils.JSONErr{Code: http.StatusInternalServerError, Msg: "unknown error"})
			if err != nil {
				logrus.Warnf("could not marsah request headers into JSON - tracing data maybe not be as expected")
			} else {
				sp.SetTag("http.response_body", string(errStatus))
			}

			// write error code to response
			w.WriteHeader(http.StatusInternalServerError)

			// write error message to response
			if err := json.NewEncoder(w).Encode(utils.JSONErr{Code: http.StatusInternalServerError, Msg: "unknown error"}); err != nil {
				logrus.Fatal(err.Error())
			}

		}
	}
}
