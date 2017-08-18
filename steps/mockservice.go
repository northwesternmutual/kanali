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
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
  "strconv"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/northwesternmutual/kanali/controller"
  "github.com/northwesternmutual/kanali/metrics"
	"github.com/northwesternmutual/kanali/spec"
	"github.com/northwesternmutual/kanali/utils"
	"github.com/opentracing/opentracing-go"
)

type mock []route

type route struct {
	Route  string      `json:"route"`
	Code   int         `json:"code"`
	Method string      `json:"method"`
	Body   interface{} `json:"body"`
}

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

	// the assumption here is that if you are using a mock,
	// you don't really care about response time. Hence,
	// we wont take the same performance enhacing measures
	// as we do for an actual request. We won't worry about
	// caching but instead we'll talk to the api server each
	// and every time for the configmap data
	cm, err := c.ClientSet.Core().ConfigMaps(proxy.ObjectMeta.Namespace).Get(proxy.Spec.Mock.ConfigMapName)
	if err != nil {
		return utils.StatusError{Code: http.StatusInternalServerError, Err: fmt.Errorf("the configmap %s could not be found in the namespace %s",
			proxy.Spec.Mock.ConfigMapName,
			proxy.ObjectMeta.Namespace,
		)}
	}

	// we have our config map. Now there's no guarentee that
	// it is in the format that we require. Using a TPR for this
	// is unnecessary and so we'll just be have to pay close attention
	// to what's in the config map and return an error if necessary
	mockResponse, ok := cm.Data["response"]
	if !ok {

		return utils.StatusError{Code: http.StatusInternalServerError, Err: fmt.Errorf("the configmap %s in the namespace %s is not formated correctly. a data field named 'response' is required",
			proxy.Spec.Mock.ConfigMapName,
			proxy.ObjectMeta.Namespace,
		)}

	}

	// we'll unmarshal the data into this var
	var mok mock

	// attempt to unmarshal json
	if err := json.Unmarshal([]byte(mockResponse), &mok); err != nil {
		return utils.StatusError{Code: http.StatusInternalServerError, Err: fmt.Errorf("the configmap %s in the namespace %s is not formated correctly. a json object is required",
			proxy.Spec.Mock.ConfigMapName,
			proxy.ObjectMeta.Namespace,
		)}
	}

	// alright so have an incoming path and we have a target path, if defined.
	// figure out with the diff is and search for that in the map.
	targetPath := utils.ComputeTargetPath(proxy.Spec.Path, proxy.Spec.Target, r.URL.EscapedPath())

	// this variable will hold the index in the array
	// of the matching mock response, if any
	mockRespIndex := -1

	// now we need to iterate over every route that we have and attempt
	// to find a match
	for i, currRoute := range mok {

		if strings.Compare(currRoute.Route, targetPath) == 0 && strings.Compare(strings.ToUpper(currRoute.Method), r.Method) == 0 {
			mockRespIndex = i
			break
		}

	}

	// no mock response was found for the incoming request
	if mockRespIndex < 0 {

		return utils.StatusError{Code: http.StatusNotFound, Err: fmt.Errorf("no mock response defined for the incoming path %s and target path %s",
			r.URL.EscapedPath(),
			targetPath,
		)}

	}

	// we have a mock response that could be anything
	// we're going to require it be JSON however.
	// we're going to attempt to marshal it into jsonpath
	// and pass it along
	mockBodyData, err := json.Marshal(mok[mockRespIndex].Body)
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
		Code:      mok[mockRespIndex].Code,
		Body:      bytes.NewBuffer(mockBodyData),
		HeaderMap: upstreamHeaders,
	}

  m.Add(metrics.Metric{Name: "http_response_code", Value: strconv.Itoa(mok[mockRespIndex].Code), Index: true})

	*resp = *responseRecorder.Result()

	return nil

}
