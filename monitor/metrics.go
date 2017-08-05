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
	"sync"

	"github.com/Sirupsen/logrus"
)

type key int

const (
	// MetricsKey is a constant contextual key for request metrics
	MetricsKey key = iota
)

// Metrics holds contextual metrics for the current request
type Metrics struct {
	mutex sync.RWMutex
	m     map[string]string
}

// New creates a new metrics objects
func New() Metrics {
	return Metrics{sync.RWMutex{}, map[string]string{}}
}

// GetCtxMetric is a concurrency safe function that retreives a
// specific contextual request metric
func GetCtxMetric(ctx context.Context, key string) string {
	untypedValue := ctx.Value(MetricsKey)
	if untypedValue == nil {
		logrus.Errorf("could not find a value in the provided context at the given key %s", key)
		return ""
	}
	metrics, ok := untypedValue.(Metrics)
	if !ok {
		logrus.Errorf("expected the value in the provided context at the given key %s to be of type Metrics", key)
		return ""
	}

	metrics.mutex.RLock()
	defer metrics.mutex.RUnlock()

	m, ok := metrics.m[key]
	if !ok {
		return ""
	}
	return m
}

// AddCtxMetric is a concurrency safe function that adds a
// specific contextual request metric
func AddCtxMetric(ctx context.Context, key string, value string) context.Context {
	untypedMetrics := ctx.Value(MetricsKey)
	if untypedMetrics == nil {
		logrus.Errorf("could not find a value in the provided context at the given key %s", key)
		return ctx
	}
	metrics, ok := untypedMetrics.(Metrics)
	if !ok {
		logrus.Errorf("expected the value in the provided context at the given key %s to be of type Metrics", key)
		return ctx
	}

	metrics.mutex.Lock()
	defer metrics.mutex.Unlock()

	metrics.m[key] = value
	return ctx
}
