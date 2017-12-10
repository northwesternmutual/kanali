package tracer

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/northwesternmutual/kanali/cmd/kanali/app/options"
	"github.com/northwesternmutual/kanali/pkg/tags"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/spf13/viper"
)

func HydrateSpanFromRequest(req *http.Request, span opentracing.Span) {
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

func HydrateSpanFromResponse(res *http.Response, span opentracing.Span) {
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
