package tracer

import (
	"bytes"
	"encoding/json"
	"github.com/northwesternmutual/kanali/config"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/spf13/viper"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	errorString = "error"
)

// HydrateSpanFromRequest adds tags to the given span relating to the given HTTP request
func HydrateSpanFromRequest(req *http.Request, span opentracing.Span) {

	if req == nil {
		span.SetTag(HTTPRequest, nil)
		return
	}

	span.SetTag(HTTPRequestMethod, req.Method)
	span.SetTag(HTTPRequestURLPath, req.URL.EscapedPath())
	span.SetTag(HTTPRequestURLHost, req.Host)

	if closerOne, closerTwo, err := dupReader(req.Body); err != nil {
		span.SetTag(HTTPRequestBody, errorString)
	} else {
		buf, err := ioutil.ReadAll(closerOne)
		if err != nil {
			span.SetTag(HTTPRequestBody, errorString)
		} else {
			span.SetTag(HTTPRequestBody, string(buf))
		}
		req.Body = closerTwo
	}

	jsonHeaders, err := json.Marshal(omitHeaderValues(
		req.Header,
		viper.GetString(config.FlagProxyHeaderMaskValue.GetLong()),
		viper.GetStringSlice(config.FlagProxyMaskHeaderKeys.GetLong())...,
	))
	if err != nil {
		span.SetTag(HTTPRequestHeaders, errorString)
	}
	span.SetTag(HTTPRequestHeaders, string(jsonHeaders))

	jsonQuery, err := json.Marshal(req.URL.Query())
	if err != nil {
		span.SetTag(HTTPRequestURLQuery, errorString)
	}
	span.SetTag(HTTPRequestURLQuery, string(jsonQuery))

}

// HydrateSpanFromResponse adds tags to the given span relating to the given HTTP response
func HydrateSpanFromResponse(res *http.Response, span opentracing.Span) {

	if res == nil {
		span.SetTag(HTTPResponse, nil)
		return
	}

	if closerOne, closerTwo, err := dupReader(res.Body); err != nil {
		span.SetTag(HTTPResponseBody, errorString)
	} else {
		buf, err := ioutil.ReadAll(closerOne)
		if err != nil {
			span.SetTag(HTTPResponseBody, errorString)
		} else {
			span.SetTag(HTTPResponseBody, string(buf))
		}
		res.Body = closerTwo
	}

	jsonHeaders, err := json.Marshal(omitHeaderValues(
		res.Header,
		viper.GetString(config.FlagProxyHeaderMaskValue.GetLong()),
		viper.GetStringSlice(config.FlagProxyMaskHeaderKeys.GetLong())...,
	))
	if err != nil {
		span.SetTag(HTTPResponseHeaders, errorString)
	}
	span.SetTag(HTTPResponseHeaders, string(jsonHeaders))

	span.SetTag(HTTPResponseStatusCode, res.StatusCode)

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
