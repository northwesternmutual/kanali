package tracer

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/northwesternmutual/kanali/cmd/kanali/app/options"
	"github.com/northwesternmutual/kanali/pkg/logging"
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
	logging.WithContext(nil).Error(msg)
}

func (l customLogger) Infof(msg string, args ...interface{}) {
	logging.WithContext(nil).Info(fmt.Sprintf(msg, args...))
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
	logger := logging.WithContext(nil)
	opentracing.SetGlobalTracer(t.tracer)

	<-ctx.Done()

	logger.Debug("tracing controller will begin gracefull termination")
	if err := t.closer.Close(); err != nil {
		logger.Error("tracing controller gracefull termination failed" + err.Error())
	}
	logger.Info("tracing controller gracefull termination successful")
}
