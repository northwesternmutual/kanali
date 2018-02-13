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
	"fmt"
	"io"
	"time"

	"github.com/northwesternmutual/kanali/cmd/kanali/app/options"
	"github.com/northwesternmutual/kanali/pkg/log"
	"github.com/opentracing/opentracing-go"
	"github.com/spf13/viper"
	jaegerConfig "github.com/uber/jaeger-client-go/config"
)

type tracerConfig struct {
	tracer opentracing.Tracer
	closer io.Closer
}

type customLogger struct{}

func (l customLogger) Error(msg string) {
	log.WithContext(nil).Error(msg)
}

func (l customLogger) Infof(msg string, args ...interface{}) {
	log.WithContext(nil).Info(fmt.Sprintf(msg, args...))
}

func Jaeger() (*tracerConfig, error) {
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

	tracer, closer, err := cfg.New("kanali", jaegerConfig.Logger(customLogger{}))
	if err != nil {
		return nil, err
	}

	return &tracerConfig{
		tracer: tracer,
		closer: closer,
	}, nil
}

func (t *tracerConfig) Run(ctx context.Context) {
	logger := log.WithContext(nil)
	opentracing.SetGlobalTracer(t.tracer)

	<-ctx.Done()

	if err := t.closer.Close(); err != nil {
		logger.Error("tracing controller gracefull termination failed" + err.Error())
	}
	logger.Info("tracing controller gracefull termination successful")
}
