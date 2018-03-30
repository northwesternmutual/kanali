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
	"net/http"
	"net/url"
	"testing"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/mocktracer"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"

	"github.com/northwesternmutual/kanali/pkg/apis/kanali.io/v2"
	"github.com/northwesternmutual/kanali/pkg/log"
	store "github.com/northwesternmutual/kanali/pkg/store/kanali/v2"
	"github.com/northwesternmutual/kanali/pkg/tags"
	"github.com/northwesternmutual/kanali/test/builder"
)

func TestOnRequest(t *testing.T) {
	core, _ := observer.New(zap.NewAtomicLevelAt(zapcore.DebugLevel))
	defer log.SetLogger(zap.New(core)).Restore()
	defer store.ApiProxyStore().Clear()

	assert.Nil(t, Plugin.OnRequest(context.Background(), nil, nil, &http.Request{
		Method: "OPTIONS",
	}))
	assert.Nil(t, Plugin.OnRequest(context.Background(), nil, nil, &http.Request{
		Method: "options",
	}))

	apiproxy := builder.NewApiProxy("foo", "bar").WithSourcePath("/foo").WithTargetBackendEndpoint("http://foo.bar.com").NewOrDie()
	u, _ := url.Parse("/foo/bar")
	assert.NotNil(t, Plugin.OnRequest(context.Background(), nil, nil, &http.Request{
		URL: u,
	}))

	store.ApiProxyStore().Set(apiproxy)
	assert.NotNil(t, Plugin.OnRequest(context.Background(), nil, nil, &http.Request{
		URL: u,
	}))

	headers := make(http.Header)
	headers.Add("apikey", "foo")
	assert.NotNil(t, Plugin.OnRequest(context.Background(), nil, nil, &http.Request{
		Header: headers,
		URL:    u,
	}))

	apikey := builder.NewApiKey("foo").WithRevision(v2.RevisionStatusActive, []byte("")).NewOrDie()
	apikey.Spec.Revisions[0].Data = "foo"
	store.ApiKeyStore().Set(apikey)

	assert.NotNil(t, Plugin.OnRequest(context.Background(), nil, nil, &http.Request{
		Header: headers,
		URL:    u,
	}))

	assert.NotNil(t, Plugin.OnRequest(context.Background(), map[string]string{
		"apiKeyBindingName": "foo",
	}, nil, &http.Request{
		Header: headers,
		URL:    u,
	}))

	store.ApiKeyBindingStore().Set(builder.NewApiKeyBinding("foo", "bar").WithKeys(
		builder.NewKeyAccess("bar").WithDefaultRule(
			builder.NewRule().WithGlobal().NewOrDie(),
		).NewOrDie(),
	).NewOrDie())

	assert.NotNil(t, Plugin.OnRequest(context.Background(), map[string]string{
		"apiKeyBindingName": "foo",
	}, nil, &http.Request{
		Header: headers,
		URL:    u,
	}))

	store.ApiKeyBindingStore().Set(builder.NewApiKeyBinding("foo", "bar").WithKeys(
		builder.NewKeyAccess("foo").WithDefaultRule(
			builder.NewRule().WithGlobal().NewOrDie(),
		).NewOrDie(),
	).NewOrDie())

	sp := mocktracer.New().StartSpan("mock").(*mocktracer.MockSpan)
	assert.Nil(t, Plugin.OnRequest(opentracing.ContextWithSpan(context.Background(), sp), map[string]string{
		"apiKeyBindingName": "foo",
	}, nil, &http.Request{
		Header: headers,
		URL:    u,
	}))
	assert.Equal(t, "foo", sp.Tag(tags.KanaliApiKeyName))
	assert.Equal(t, "foo", sp.Tag(tags.KanaliApiKeyBindingName))
	assert.Equal(t, "bar", sp.Tag(tags.KanaliApiKeyBindingNamespace))

	store.ApiKeyBindingStore().Set(builder.NewApiKeyBinding("foo", "bar").WithKeys(
		builder.NewKeyAccess("foo").WithDefaultRule(
			builder.NewRule().NewOrDie(),
		).NewOrDie(),
	).NewOrDie())

	assert.NotNil(t, Plugin.OnRequest(context.Background(), map[string]string{
		"apiKeyBindingName": "foo",
	}, nil, &http.Request{
		Header: headers,
		URL:    u,
	}))
}

func TestOnResponse(t *testing.T) {
	assert.Nil(t, Plugin.OnResponse(context.Background(), nil, nil, nil))
}

func TestValidateApiKey(t *testing.T) {
	assert.False(t, validateApiKey(nil, ""))

	assert.True(t, validateApiKey(&v2.Rule{
		Global: true,
		Granular: v2.GranularProxy{
			Verbs: []string{},
		},
	}, "GET"), "rule should be authorized")

	assert.True(t, validateApiKey(&v2.Rule{
		Global: true,
		Granular: v2.GranularProxy{
			Verbs: []string{
				"GET",
			},
		},
	}, "GET"), "rule should be authorized")

	assert.True(t, validateApiKey(&v2.Rule{
		Global: false,
		Granular: v2.GranularProxy{
			Verbs: []string{
				"GET",
			},
		},
	}, "GET"), "rule should be authorized")

	assert.True(t, validateApiKey(&v2.Rule{
		Global: true,
		Granular: v2.GranularProxy{
			Verbs: []string{
				"POST",
			},
		},
	}, "GET"), "rule should be authorized")

	assert.False(t, validateApiKey(&v2.Rule{
		Global: false,
		Granular: v2.GranularProxy{
			Verbs: []string{
				"POST",
			},
		},
	}, "GET"), "rule should not be authorized")
}

func TestValidateGranularRules(t *testing.T) {
	assert.False(t, validateGranularRules("GET", v2.GranularProxy{
		Verbs: []string{},
	}))

	assert.True(t, validateGranularRules("GET", v2.GranularProxy{
		Verbs: []string{
			"GET",
		},
	}), "http method should be authorized")

	assert.True(t, validateGranularRules("get", v2.GranularProxy{
		Verbs: []string{
			"GET",
			"POST",
			"PUT",
		},
	}), "http method should be authorized")

	assert.True(t, validateGranularRules("put", v2.GranularProxy{
		Verbs: []string{
			"GET",
			"POST",
			"PUT",
		},
	}), "http method should be authorized")

	assert.False(t, validateGranularRules("GET", v2.GranularProxy{
		Verbs: []string{
			"POST",
		},
	}), "http method should be authorized")

	assert.False(t, validateGranularRules("get", v2.GranularProxy{
		Verbs: []string{
			"POST",
			"PUT",
		},
	}), "http method should be authorized")

	assert.False(t, validateGranularRules("", v2.GranularProxy{
		Verbs: []string{
			"POST",
		},
	}), "http method should be authorized")

	assert.False(t, validateGranularRules("HTTP", v2.GranularProxy{
		Verbs: []string{
			"POST",
		},
	}), "http method should be authorized")
}

func TestExtractApiKey(t *testing.T) {
	headers := make(http.Header)
	headers.Add("apikey", "foo")

	result, _ := extractApiKey(headers)
	assert.Equal(t, "foo", result)

	_, err := extractApiKey(nil)
	assert.NotNil(t, err)

	headers.Set("apikey", "")
	_, err = extractApiKey(headers)
	assert.NotNil(t, err)
}
