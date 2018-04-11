// Copyright (c) 2018 Northwestern Mutual.
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

package validate

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"

	"github.com/northwesternmutual/kanali/pkg/apis/kanali.io/v2"
	"github.com/northwesternmutual/kanali/pkg/client/clientset/versioned/fake"
	"github.com/northwesternmutual/kanali/pkg/log"
	"github.com/northwesternmutual/kanali/test/builder"
)

func TestIsValidApiProxy(t *testing.T) {
	core, _ := observer.New(zap.NewAtomicLevelAt(zapcore.DebugLevel))
	defer log.SetLogger(zap.New(core)).Restore()

	assert.Error(t, (&validation{}).IsValidApiProxy(nil))

	tests := []struct {
		isValid           bool
		testApiProxy      *v2.ApiProxy
		presentApiProxies []*v2.ApiProxy
	}{
		{
			isValid:           true,
			testApiProxy:      builder.NewApiProxy("foo", "bar").WithSourcePath("/foo").NewOrDie(),
			presentApiProxies: nil,
		},
		{
			isValid:      true,
			testApiProxy: builder.NewApiProxy("foo", "bar").WithSourcePath("/foo").NewOrDie(),
			presentApiProxies: []*v2.ApiProxy{
				builder.NewApiProxy("bar", "foo").WithSourcePath("/bar").NewOrDie(),
			},
		},
		{
			isValid:      true,
			testApiProxy: builder.NewApiProxy("foo", "bar").WithSourcePath("/foo").NewOrDie(),
			presentApiProxies: []*v2.ApiProxy{
				builder.NewApiProxy("foo", "bar").WithSourcePath("/foo").NewOrDie(),
			},
		},
		{
			isValid:      true,
			testApiProxy: builder.NewApiProxy("foo", "bar").WithSourcePath("/foo").WithSourceHost("foo.bar.com").NewOrDie(),
			presentApiProxies: []*v2.ApiProxy{
				builder.NewApiProxy("foo", "bar").WithSourcePath("/foo").NewOrDie(),
			},
		},
		{
			isValid:      false,
			testApiProxy: builder.NewApiProxy("bar", "foo").WithSourcePath("/foo").NewOrDie(),
			presentApiProxies: []*v2.ApiProxy{
				builder.NewApiProxy("foo", "bar").WithSourcePath("/foo").NewOrDie(),
			},
		},
		{
			isValid:      false,
			testApiProxy: builder.NewApiProxy("bar", "foo").WithSourcePath("/foo").WithSourceHost("foo.bar.com").NewOrDie(),
			presentApiProxies: []*v2.ApiProxy{
				builder.NewApiProxy("foo", "bar").WithSourcePath("/foo").WithSourceHost("foo.bar.com").NewOrDie(),
			},
		},
		{
			isValid:      true,
			testApiProxy: builder.NewApiProxy("bar", "foo").WithSourcePath("/foo").WithSourceHost("foo.bar.com").NewOrDie(),
			presentApiProxies: []*v2.ApiProxy{
				builder.NewApiProxy("bar", "foo").WithSourcePath("/bar").WithSourceHost("foo.bar.com").NewOrDie(),
			},
		},
		{
			isValid:      true,
			testApiProxy: builder.NewApiProxy("foo", "bar").WithSourcePath("/foo").NewOrDie(),
			presentApiProxies: []*v2.ApiProxy{
				builder.NewApiProxy("foo", "bar").WithSourcePath("/foo").WithSourceHost("foo.bar.com").NewOrDie(),
			},
		},
		{
			isValid:      false,
			testApiProxy: builder.NewApiProxy("foo", "bar").WithSourcePath("/foo").WithSourceHost("foo.bar.com").NewOrDie(),
			presentApiProxies: []*v2.ApiProxy{
				builder.NewApiProxy("bar", "foo").WithSourcePath("/foo").WithSourceHost("foo.bar.com").NewOrDie(),
			},
		},
		{
			isValid:      true,
			testApiProxy: builder.NewApiProxy("foo", "bar").WithSourcePath("/foo").WithSourceHost("bar.foo.com").NewOrDie(),
			presentApiProxies: []*v2.ApiProxy{
				builder.NewApiProxy("foo", "bar").WithSourcePath("/foo").WithSourceHost("foo.bar.com").NewOrDie(),
			},
		},
		{
			isValid:      false,
			testApiProxy: builder.NewApiProxy("foo", "bar").WithSourcePath("/foo").NewOrDie(),
			presentApiProxies: []*v2.ApiProxy{
				builder.NewApiProxy("bar", "foo").WithSourcePath("/foo").WithSourceHost("foo.bar.com").NewOrDie(),
			},
		},
		{
			isValid:      true,
			testApiProxy: builder.NewApiProxy("foo", "bar").WithSourcePath("/foo").WithSourceHost("foo.bar.com").NewOrDie(),
			presentApiProxies: []*v2.ApiProxy{
				builder.NewApiProxy("foo", "bar").WithSourcePath("/foo").NewOrDie(),
			},
		},
		{
			isValid:      false,
			testApiProxy: builder.NewApiProxy("foo", "bar").WithSourcePath("/foo").NewOrDie(),
			presentApiProxies: []*v2.ApiProxy{
				builder.NewApiProxy("bar", "foo").WithSourcePath("/foo").WithSourceHost("bar.foo.com").NewOrDie(),
				builder.NewApiProxy("car", "baz").WithSourcePath("/foo").WithSourceHost("foo.bar.com").NewOrDie(),
			},
		},
		{
			isValid:      false,
			testApiProxy: builder.NewApiProxy("bar", "foo").WithSourcePath("/foo").NewOrDie(),
			presentApiProxies: []*v2.ApiProxy{
				builder.NewApiProxy("bar", "foo").WithSourcePath("/foo").WithSourceHost("bar.foo.com").NewOrDie(),
				builder.NewApiProxy("car", "baz").WithSourcePath("/foo").WithSourceHost("foo.bar.com").NewOrDie(),
			},
		},
		{
			isValid:      false,
			testApiProxy: builder.NewApiProxy("bar", "foo").WithSourcePath("/foo").WithSourceHost("foo.bar.com").NewOrDie(),
			presentApiProxies: []*v2.ApiProxy{
				builder.NewApiProxy("bar", "foo").WithSourcePath("/foo").WithSourceHost("bar.foo.com").NewOrDie(),
				builder.NewApiProxy("car", "baz").WithSourcePath("/foo").WithSourceHost("foo.bar.com").NewOrDie(),
			},
		},
		{
			isValid:      false,
			testApiProxy: builder.NewApiProxy("foo", "bar").WithSourcePath("/foo").NewOrDie(),
			presentApiProxies: []*v2.ApiProxy{
				builder.NewApiProxy("foo", "bar").WithSourcePath("/foo").WithSourceHost("bar.foo.com").NewOrDie(),
				builder.NewApiProxy("bar", "foo").WithSourcePath("/foo").WithSourceHost("foo.bar.com").NewOrDie(),
			},
		},
		{
			isValid:      true,
			testApiProxy: builder.NewApiProxy("bar", "foo").WithSourcePath("/foo").WithSourceHost("foo.bar.com").NewOrDie(),
			presentApiProxies: []*v2.ApiProxy{
				builder.NewApiProxy("bar", "foo").WithSourcePath("/foo").WithSourceHost("bar.foo.com").NewOrDie(),
				builder.NewApiProxy("car", "baz").WithSourcePath("/foo").WithSourceHost("car.baz.com").NewOrDie(),
			},
		},
	}

	for _, test := range tests {
		validator := New(context.Background(), fake.NewSimpleClientset())
		for _, apiproxy := range test.presentApiProxies {
			_, err := validator.clientset.KanaliV2().ApiProxies(apiproxy.GetNamespace()).Create(apiproxy)
			assert.Nil(t, err)
		}
		apiproxyData, err := json.Marshal(test.testApiProxy)
		assert.Nil(t, err)

		result := validator.IsValidApiProxy(apiproxyData)
		if test.isValid {
			assert.Nil(t, result)
		} else {
			assert.NotNil(t, result)
		}
	}
}

func TestIsValidApiProxyList(t *testing.T) {
	core, _ := observer.New(zap.NewAtomicLevelAt(zapcore.DebugLevel))
	defer log.SetLogger(zap.New(core)).Restore()

	validator := New(context.Background(), fake.NewSimpleClientset())
	assert.Error(t, (&validation{}).IsValidApiProxyList(nil))

	apiproxyListData, err := json.Marshal(&v2.ApiProxyList{
		Items: []v2.ApiProxy{
			*builder.NewApiProxy("foo", "bar").WithSourcePath("/foo").WithSourceHost("foo.bar.com").NewOrDie(),
		},
	})
	assert.Nil(t, err)
	assert.Nil(t, validator.IsValidApiProxyList(apiproxyListData))

	_, err = validator.clientset.KanaliV2().ApiProxies("foo").Create(
		builder.NewApiProxy("bar", "foo").WithSourcePath("/foo").NewOrDie(),
	)
	assert.Nil(t, err)
	assert.NotNil(t, validator.IsValidApiProxyList(apiproxyListData))
}
