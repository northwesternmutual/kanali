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
  "fmt"
	"bytes"
	"context"
	"encoding/json"
	"errors"
  "strconv"
	"net/http"
	"net/http/httptest"

	"github.com/northwesternmutual/kanali/controller"
  "github.com/northwesternmutual/kanali/metrics"
	"github.com/northwesternmutual/kanali/spec"
	"github.com/northwesternmutual/kanali/utils"
	"github.com/opentracing/opentracing-go"
)

// MockServiceStep is factory that defines a step responsible for
// discovering a mock response for the incoming request
type MockServiceStep struct{}

// GetName retruns the name of the MockServiceStep step
func (step MockServiceStep) GetName() string {
	return "Mock Service"
}

// Do executes the logic of the MockServiceStep step
func (step MockServiceStep) Do(ctx context.Context, m *metrics.Metrics, c *controller.Controller, w http.ResponseWriter, r *http.Request, resp *http.Response, trace opentracing.Span) error {

	// incoming proxy
	untypedProxy, err := spec.ProxyStore.Get(r.URL.EscapedPath())
	if err != nil || untypedProxy == nil {
		return utils.StatusError{Code: http.StatusNotFound, Err: errors.New("proxy not found")}
	}

	proxy, ok := untypedProxy.(spec.APIProxy)
	if !ok {
		return utils.StatusError{Code: http.StatusNotFound, Err: errors.New("proxy not found")}
	}

  targetPath := utils.ComputeTargetPath(proxy.Spec.Path, proxy.Spec.Target, r.URL.EscapedPath())

  untypedMr, err := spec.MockResponseStore.Get(proxy.ObjectMeta.Namespace, proxy.Spec.Mock.ConfigMapName, targetPath, r.Method)
  if err != nil {
    return utils.StatusError{Code: http.StatusInternalServerError, Err: fmt.Errorf("error retrieving mock response: %s", err.Error())}
  }
  if untypedMr == nil {
    return utils.StatusError{Code: http.StatusNotFound, Err: errors.New("no mock response found")}
  }
  mr, ok := untypedMr.(spec.Route)
	if !ok {
		return utils.StatusError{Code: http.StatusNotFound, Err: errors.New("no mock response found")}
	}

	mockBodyData, err := json.Marshal(mr.Body)
	if err != nil {
		return utils.StatusError{Code: http.StatusInternalServerError, Err: fmt.Errorf("the configmap %s in the namespace %s is not formated correctly. while data was found for the incoming route, it was not valid json",
			proxy.Spec.Mock.ConfigMapName,
			proxy.ObjectMeta.Namespace,
		)}
	}

	// create new upstream header object
	upstreamHeaders := http.Header{}

	// currently, we are enforcing a json response
	upstreamHeaders.Add("Content-Type", "application/json")

	// create a fake response
	responseRecorder := &httptest.ResponseRecorder{
		Code:      mr.Code,
		Body:      bytes.NewBuffer(mockBodyData),
		HeaderMap: upstreamHeaders,
	}

  m.Add(metrics.Metric{Name: "http_response_code", Value: strconv.Itoa(mr.Code), Index: true})

	*resp = *responseRecorder.Result()

	return nil

}
