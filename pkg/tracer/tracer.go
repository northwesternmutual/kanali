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

package tracer

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/ghodss/yaml"

	"github.com/northwesternmutual/kanali/cmd/kanali/app/options"
	"github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/viper"
	jaegerConfig "github.com/uber/jaeger-client-go/config"
	jaegerProm "github.com/uber/jaeger-lib/metrics/prometheus"

	"github.com/northwesternmutual/kanali/pkg/log"
)

type tracerParams struct {
	tracer opentracing.Tracer
	closer io.Closer
}

func New() (*tracerParams, error) {

	cfg, err := parseConfig(viper.GetString(options.FlagTracingConfig.GetLong()))
	if err != nil {
		return nil, err
	}

	tracer, closer, err := (*cfg).New("kanali",
		jaegerConfig.Gen128Bit(true),
		jaegerConfig.Logger(customLogger{}),
		jaegerConfig.Metrics(jaegerProm.New(
			jaegerProm.WithRegisterer(prometheus.DefaultRegisterer),
		)),
	)
	if err != nil {
		return nil, err
	}

	return &tracerParams{
		tracer: tracer,
		closer: closer,
	}, nil
}

func parseConfig(location string) (*jaegerConfig.Configuration, error) {
	f, err := os.Open(location)
	if err != nil {
		return nil, err
	} else {
		defer func() {
			if err := f.Close(); err != nil {
				log.WithContext(nil).Error(fmt.Sprintf("error closing tracer: %s", err))
			}
		}()
	}
	return doParseConfig(f)
}

func doParseConfig(r io.Reader) (*jaegerConfig.Configuration, error) {
	if r == nil {
		return nil, errors.New("reader is nil")
	}

	var cfg jaegerConfig.Configuration
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (params *tracerParams) Run(ctx context.Context) error {
	opentracing.SetGlobalTracer(params.tracer)
	<-ctx.Done()
	return nil
}

func (params *tracerParams) Close(error) error {
	return params.closer.Close()
}
