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

package apikey

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	//"github.com/northwesternmutual/kanali/pkg/metrics"
	"github.com/northwesternmutual/kanali/pkg/apis/kanali.io/v2"
	kanaliErrors "github.com/northwesternmutual/kanali/pkg/errors"
	"github.com/northwesternmutual/kanali/pkg/log"
	store "github.com/northwesternmutual/kanali/pkg/store/kanali/v2"
	"github.com/northwesternmutual/kanali/pkg/tags"
	"github.com/northwesternmutual/kanali/pkg/traffic"
	"github.com/northwesternmutual/kanali/pkg/utils"
	pluginConfig "github.com/northwesternmutual/kanali/plugins/apikey/config"
	opentracing "github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

// ApiKeyFactory is factory that implements the github.com/northwesternmutual/kanali/pkg/plugin.Plugin interface
type ApiKeyFactory struct{}

// OnRequest intercepts a request before it get proxied to an upstream service
func (k ApiKeyFactory) OnRequest(ctx context.Context, config map[string]string, w *httptest.ResponseRecorder, r *http.Request) error {

	trafficCtlr, _ := traffic.NewController()

	logger := log.WithContext(r.Context())

	p := store.ApiProxyStore().Get(utils.ComputeURLPath(r.URL))
	if p == nil {
		logger.Warn(kanaliErrors.ErrorProxyNotFound.Message)
		return kanaliErrors.ErrorProxyNotFound
	}

	span := opentracing.SpanFromContext(ctx)

	timestamp := time.Now()

	// do not continue if an OPTION request
	if strings.ToUpper(r.Method) == "OPTIONS" {
		logger.Debug("api key validation is not preformed for OPTIONS requests")
		return next()
	}

	// extract the api key
	apiKeyText, err := extractApiKey(r.Header)
	if err != nil {
		logger.Error(err.Error())
		return kanaliErrors.ErrorForbidden
	}

	// attempt to find a matching api key
	apiKeyObj := store.ApiKeyStore().Get([]byte(apiKeyText))
	if apiKeyObj == nil {
		logger.Warn("api key was not found in store")
		return kanaliErrors.ErrorApiKeyUnauthorized
	}

	cfg, err := pluginConfig.New(config)
	if err != nil {
		logger.Error(err.Error())
		return kanaliErrors.ErrorApiKeyUnauthorized
	}

	// BEGIN logging, metrics, and tracing overhead
	logger.Debug("ApiKey resource details",
		zap.String(tags.KanaliApiKeyName, apiKeyObj.ObjectMeta.Name),
	)
	span.SetTag(tags.KanaliApiKeyName, apiKeyObj.ObjectMeta.Name)
	// m.Add(metrics.Metric{
	// 	Name:  tags.KanaliApiKeyName,
	// 	Value: apiKeyObj.ObjectMeta.Name,
	// 	Index: true,
	// })
	// END logging, metrics, and tracing overhead

	if !store.ApiKeyBindingStore().Contains(p.ObjectMeta.Namespace, cfg.ApiKeyBindingName) {
		logger.Warn("ApiKeyBinding resource was not found in store",
			zap.String(tags.KanaliApiKeyBindingName, cfg.ApiKeyBindingName),
			zap.String(tags.KanaliApiKeyBindingNamespace, p.ObjectMeta.Namespace),
		)
		return kanaliErrors.ErrorApiKeyUnauthorized
	}

	span.SetTag(tags.KanaliApiKeyBindingName, cfg.ApiKeyBindingName)
	span.SetTag(tags.KanaliApiKeyBindingNamespace, p.ObjectMeta.Namespace)
	logger.Info("ApiKey resource details",
		zap.String(tags.KanaliApiKeyBindingName, cfg.ApiKeyBindingName),
		zap.String(tags.KanaliApiKeyBindingNamespace, p.ObjectMeta.Namespace),
	)

	if !store.ApiKeyBindingStore().ContainsApiKey(p.ObjectMeta.Namespace, cfg.ApiKeyBindingName, apiKeyObj.ObjectMeta.Name) {
		logger.Error("ApiKeyBinding resource does not any permissions for given ApiKey resource",
			zap.String(tags.KanaliApiKeyBindingName, cfg.ApiKeyBindingName),
			zap.String(tags.KanaliApiKeyBindingNamespace, p.ObjectMeta.Namespace),
			zap.String(tags.KanaliApiKeyName, apiKeyObj.ObjectMeta.Name),
		)
		return kanaliErrors.ErrorApiKeyUnauthorized
	}

	rule, rate := store.ApiKeyBindingStore().GetRuleAndRate(p.ObjectMeta.Namespace, cfg.ApiKeyBindingName, apiKeyObj.ObjectMeta.Name, utils.ComputeTargetPath(p.Spec.Source.Path, p.Spec.Target.Path, r.URL.Path))

	if !validateApiKey(rule, r.Method) {
		return kanaliErrors.ErrorApiKeyUnauthorized
	}

	if store.TrafficStore().IsRateLimitViolated(*p, rate, apiKeyObj.ObjectMeta.Name, timestamp) {
		logger.Info("rate limit exceeded")
		return kanaliErrors.ErrorTooManyRequests
	}

	go trafficCtlr.Report(ctx, &store.TrafficPoint{
		Time:      timestamp.UnixNano(),
		Namespace: p.ObjectMeta.Namespace,
		ProxyName: config["apiKeyBindingName"],
		KeyName:   apiKeyObj.ObjectMeta.Name,
	})

	return next()

}

// OnResponse intercepts a request after it has been proxied to an upstream service
// but before the response gets returned to the client
func (k ApiKeyFactory) OnResponse(ctx context.Context, config map[string]string, w *httptest.ResponseRecorder, r *http.Request) error {
	log.WithContext(ctx).Info("api key plugin OnRequest method not implemented")
	return next()
}

// validateApiKey will return true if the given api key
// is authorized to make the given request.
// Global rule valudation will be given priority over
// granular rule validation
func validateApiKey(rule *v2.Rule, method string) bool {
	if rule == nil {
		return false
	}
	return rule.Global || validateGranularRules(method, rule.Granular)
}

// check to see wheather a given HTTP method can be found
// in the list of HTTP methods belonging to a spec.GranularProxy
func validateGranularRules(method string, rule v2.GranularProxy) bool {
	if len(rule.Verbs) < 1 {
		return false
	}
	for _, verb := range rule.Verbs {
		if strings.ToUpper(verb) == strings.ToUpper(method) {
			return true
		}
	}
	return false
}

func next() error {
	return nil
}

func extractApiKey(reqHeaders http.Header) (string, error) {
	apiKeyText := reqHeaders.Get("apikey")
	if len(apiKeyText) < 1 {
		return "", errors.New("expected the apikey header to contain an api key value")
	}
	return apiKeyText, nil
}

// Plugin can be discovered by golang plugin package
var Plugin ApiKeyFactory
