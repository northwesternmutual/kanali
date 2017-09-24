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

	"github.com/northwesternmutual/kanali/config"
	"github.com/northwesternmutual/kanali/flow"
	"github.com/northwesternmutual/kanali/metrics"
	"github.com/northwesternmutual/kanali/spec"
	"github.com/northwesternmutual/kanali/steps"
	"github.com/opentracing/opentracing-go"
	"github.com/spf13/viper"
)

// IncomingRequest orchestrates the logic that occurs for every incoming request
func IncomingRequest(ctx context.Context, proxy *spec.APIProxy, m *metrics.Metrics, w http.ResponseWriter, r *http.Request, trace opentracing.Span) error {

	// this is a handler to our future proxy pass response
	// maybe there's a better way to do this... seems misplaced
	futureResponse := &http.Response{}

	f := &flow.Flow{}

	f.Add(
		steps.ValidateProxyStep{},
		steps.PluginsOnRequestStep{},
	)
	if viper.GetBool(config.FlagProxyEnableMockResponses.GetLong()) && mockIsDefined(r.URL.Path) {
		f.Add(steps.MockServiceStep{})
	} else {
		f.Add(steps.ProxyPassStep{})
	}

	f.Add(
		steps.PluginsOnResponseStep{},
		steps.WriteResponseStep{},
	)

	err := f.Play(ctx, proxy, m, w, r, futureResponse, trace)

	return err

}

func mockIsDefined(path string) bool {

	untypedProxy, err := spec.ProxyStore.Get(path)
	if err != nil || untypedProxy == nil {
		return false
	}

	proxy, ok := untypedProxy.(spec.APIProxy)
	if !ok {
		return false
	}

	if proxy.Spec.Mock != nil {
		return proxy.Spec.Mock.ConfigMapName != ""
	}

	return false

}
