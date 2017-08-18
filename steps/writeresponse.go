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

package steps

import (
	"context"
	"io"
	"net/http"
	"strconv"

	"github.com/Sirupsen/logrus"
	"github.com/northwesternmutual/kanali/controller"
	"github.com/northwesternmutual/kanali/metrics"
	"github.com/northwesternmutual/kanali/utils"
	"github.com/opentracing/opentracing-go"
)

// WriteResponseStep is factory that defines a step responsible for writing
// an HTTP response
type WriteResponseStep struct{}

// GetName retruns the name of the WriteResponseStep step
func (step WriteResponseStep) GetName() string {
	return "Write Response"
}

// Do executes the logic of the WriteResponseStep step
func (step WriteResponseStep) Do(ctx context.Context, m *metrics.Metrics, c *controller.Controller, w http.ResponseWriter, r *http.Request, resp *http.Response, trace opentracing.Span) error {

	for k, v := range resp.Header {
		for _, value := range v {
			w.Header().Set(k, value)
		}
	}

	m.Add(metrics.Metric{Name: "http_response_code", Value: strconv.Itoa(resp.StatusCode), Index: true})

	closer, str, err := utils.DupReaderAndString(resp.Body)
	if err != nil {
		logrus.Errorf("error copying response body, response may not be as expected: %s", err.Error())
	}

	trace.SetTag("http.status_code", resp.StatusCode)
	trace.SetTag("http.response_body", str)

	w.WriteHeader(resp.StatusCode)

	if _, err := io.Copy(w, closer); err != nil {
		return err
	}

	return nil
}
