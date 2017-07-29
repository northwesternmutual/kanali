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

package monitor

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddCtxMetric(t *testing.T) {
	assert := assert.New(t)

	ctx := context.WithValue(context.Background(), MetricsKey, New())
	ctx = AddCtxMetric(ctx, "foo", "bar")
	ctx = AddCtxMetric(ctx, "one", "two")

	untypedValue := ctx.Value(MetricsKey)
	if untypedValue == nil {
		assert.Fail("context value should not be nil")
	}

	value, ok := untypedValue.(Metrics)
	if !ok {
		assert.Fail("value should be of type Metrics")
	}

	assert.Equal(value.m["foo"], "bar", "map value not incorrect")
	assert.Equal(value.m["one"], "two", "map value not incorrect")

	ctx = context.Background()
	ctx = AddCtxMetric(ctx, "foo", "bar")
	ctx = AddCtxMetric(ctx, "one", "two")

	untypedValue = ctx.Value(MetricsKey)
	if untypedValue == nil {
		assert.Fail("context value should not be nil")
	}

	value, ok = untypedValue.(Metrics)
	if !ok {
		assert.Fail("value should be of type Metrics")
	}

	assert.Equal(value.m["foo"], "bar", "map value not incorrect")
	assert.Equal(value.m["one"], "two", "map value not incorrect")
}

func TestGetCtxMetric(t *testing.T) {

	assert := assert.New(t)

	val := GetCtxMetric(context.Background(), "foo")
	assert.Equal(val, "")

	ctx := context.WithValue(context.Background(), MetricsKey, "foo")
	val = GetCtxMetric(ctx, "foo")
	assert.Equal(val, "")

	ctx = context.WithValue(context.Background(), MetricsKey, New())
	ctx = AddCtxMetric(ctx, "foo", "bar")
	ctx = AddCtxMetric(ctx, "one", "two")

	val = GetCtxMetric(ctx, "foo")
	assert.Equal(val, "bar", "wrong value from metric map")
	val = GetCtxMetric(ctx, "one")
	assert.Equal(val, "two", "wrong value from metric map")
	val = GetCtxMetric(ctx, "frank")
	assert.Equal(val, "", "value should not exist")

}
