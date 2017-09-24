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

package tracer

import (
	"fmt"
	"io"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/northwesternmutual/kanali/config"
	"github.com/opentracing/opentracing-go"
	"github.com/spf13/viper"
	jaegerConfig "github.com/uber/jaeger-client-go/config"
)

type customLogger struct{}

func (l customLogger) Error(msg string) {
	logrus.Error(msg)
}

func (l customLogger) Infof(msg string, args ...interface{}) {
	logrus.Info(fmt.Sprintf(msg, args...))
}

// Jaeger creates a new opentracing compatible tracer
func Jaeger() (opentracing.Tracer, io.Closer, error) {

	cfg := jaegerConfig.Configuration{
		Sampler: &jaegerConfig.SamplerConfig{
			Type:              "const",
			SamplingServerURL: viper.GetString(config.FlagTracingJaegerServerURL.GetLong()),
			Param:             1,
		},
		Reporter: &jaegerConfig.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
			LocalAgentHostPort:  fmt.Sprintf("%s:5775", viper.GetString(config.FlagTracingJaegerAgentURL.GetLong())),
		},
	}

	return cfg.New("kanali", jaegerConfig.Logger(customLogger{}))

}
