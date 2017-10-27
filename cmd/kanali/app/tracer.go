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

package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/northwesternmutual/kanali/cmd/kanali/app/options"
	"github.com/northwesternmutual/kanali/pkg/logging"
	tags "github.com/northwesternmutual/kanali/pkg/tags"
	"github.com/opentracing/opentracing-go"
	"github.com/spf13/viper"
	jaegerConfig "github.com/uber/jaeger-client-go/config"
)

type customLogger struct{}

func (l customLogger) Error(msg string) {
	logging.WithContext(nil).Error(msg)
}

func (l customLogger) Infof(msg string, args ...interface{}) {
	logging.WithContext(nil).Info(fmt.Sprintf(msg, args...))
}

func newJaegerTracer() (opentracing.Tracer, io.Closer, error) {
	cfg := jaegerConfig.Configuration{
		Sampler: &jaegerConfig.SamplerConfig{
			Type:              "const",
			SamplingServerURL: viper.GetString(options.FlagTracingJaegerServerURL.GetLong()),
			Param:             1,
		},
		Reporter: &jaegerConfig.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
			LocalAgentHostPort:  fmt.Sprintf("%s:5775", viper.GetString(options.FlagTracingJaegerAgentURL.GetLong())),
		},
	}

	return cfg.New("kanali", jaegerConfig.Logger(customLogger{}))
}

func hydrateSpanFromRequest(req *http.Request, span opentracing.Span) {
	if req == nil {
		span.SetTag(tags.HTTPRequest, nil)
		return
	}

	span.SetTag(tags.HTTPRequestMethod, req.Method)
	span.SetTag(tags.HTTPRequestURLPath, req.URL.EscapedPath())
	span.SetTag(tags.HTTPRequestURLHost, req.Host)

	if closerOne, closerTwo, err := dupReader(req.Body); err != nil {
		span.SetTag(tags.HTTPRequestBody, tags.Error)
	} else {
		buf, err := ioutil.ReadAll(closerOne)
		if err != nil {
			span.SetTag(tags.HTTPRequestBody, tags.Error)
		} else {
			span.SetTag(tags.HTTPRequestBody, string(buf))
		}
		req.Body = closerTwo
	}

	jsonHeaders, err := json.Marshal(omitHeaderValues(
		req.Header,
		viper.GetString(options.FlagProxyHeaderMaskValue.GetLong()),
		viper.GetStringSlice(options.FlagProxyMaskHeaderKeys.GetLong())...,
	))
	if err != nil {
		span.SetTag(tags.HTTPRequestHeaders, tags.Error)
	}
	span.SetTag(tags.HTTPRequestHeaders, string(jsonHeaders))

	jsonQuery, err := json.Marshal(req.URL.Query())
	if err != nil {
		span.SetTag(tags.HTTPRequestURLQuery, tags.Error)
	}
	span.SetTag(tags.HTTPRequestURLQuery, string(jsonQuery))
}

func hydrateSpanFromResponse(res *http.Response, span opentracing.Span) {
	if res == nil {
		span.SetTag(tags.HTTPResponse, nil)
		return
	}

	if closerOne, closerTwo, err := dupReader(res.Body); err != nil {
		span.SetTag(tags.HTTPResponseBody, tags.Error)
	} else {
		buf, err := ioutil.ReadAll(closerOne)
		if err != nil {
			span.SetTag(tags.HTTPResponseBody, tags.Error)
		} else {
			span.SetTag(tags.HTTPResponseBody, string(buf))
		}
		res.Body = closerTwo
	}

	jsonHeaders, err := json.Marshal(omitHeaderValues(
		res.Header,
		viper.GetString(options.FlagProxyHeaderMaskValue.GetLong()),
		viper.GetStringSlice(options.FlagProxyMaskHeaderKeys.GetLong())...,
	))
	if err != nil {
		span.SetTag(tags.HTTPResponseHeaders, tags.Error)
	}
	span.SetTag(tags.HTTPResponseHeaders, string(jsonHeaders))
	span.SetTag(tags.HTTPResponseStatusCode, res.StatusCode)
}

func dupReader(closer io.ReadCloser) (io.ReadCloser, io.ReadCloser, error) {

	buf, err := ioutil.ReadAll(closer)
	if err != nil {
		return nil, nil, err
	}

	rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))
	rdr2 := ioutil.NopCloser(bytes.NewBuffer(buf))

	return rdr1, rdr2, nil

}

func omitHeaderValues(h http.Header, msg string, keys ...string) http.Header {
	if h == nil {
		return http.Header{}
	}
	copy := http.Header{}
	for k, v := range h {
		copy[strings.Title(k)] = v
	}
	for _, key := range keys {
		if copy.Get(key) != "" {
			copy.Set(key, msg)
		}
	}
	return copy
}
