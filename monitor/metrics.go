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

	"github.com/Sirupsen/logrus"
)

type key int

const (
	// MetricsKey is a constant contextual key for request metrics
	MetricsKey key = iota
)

// Metrics holds contextual metrics for the current request
type Metrics struct {
	m map[string]string
}

// New creates a new metrics objects
func New() Metrics {
	return Metrics{map[string]string{}}
}

// GetCtxMetric retreives a specific contextual request metric
func GetCtxMetric(ctx context.Context, key string) string {

  // TODO: add read/write mutex lock

	untypedValue := ctx.Value(MetricsKey)
	if untypedValue == nil {
		logrus.Errorf("context does not have the correct key")
		return ""
	}
	value, ok := untypedValue.(Metrics)
	if !ok {
		logrus.Errorf("value must be of type Metrics")
		return ""
	}
	m, ok := value.m[key]
	if !ok {
		return ""
	}
	return m
}

// AddCtxMetric adds a specific contextual request metric
func AddCtxMetric(ctx context.Context, key string, value string) context.Context {

  // TODO: add read/write mutex lock

	untypedMetrics := ctx.Value(MetricsKey)
	if untypedMetrics == nil {
		ctx = context.WithValue(context.Background(), MetricsKey, Metrics{map[string]string{}})
	}
	metrics, ok := untypedMetrics.(Metrics)
	if !ok {
		metrics = Metrics{map[string]string{}}
		ctx = context.WithValue(ctx, MetricsKey, metrics)
	}
	metrics.m[key] = value
	return ctx
}
