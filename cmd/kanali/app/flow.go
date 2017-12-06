package app

import (
  "context"
  "net/http"

  "go.uber.org/zap"
  "k8s.io/client-go/informers/core"
  "github.com/northwesternmutual/kanali/pkg/logging"
  "github.com/northwesternmutual/kanali/pkg/metrics"
  opentracing "github.com/opentracing/opentracing-go"
  tags "github.com/northwesternmutual/kanali/pkg/tags"
  "github.com/northwesternmutual/kanali/pkg/apis/kanali.io/v2"
)

type flow []step

type step interface {
	getName() string
	do(ctx context.Context, proxy *v2.ApiProxy, k8sCoreClient core.Interface, metrics *metrics.Metrics, w http.ResponseWriter, r *http.Request, resp *http.Response, trace opentracing.Span) error
}

func (f *flow) add(steps ...step) {
	for _, s := range steps {
		*f = append(*f, s)
	}
}

func (f *flow) play(ctx context.Context, proxy *v2.ApiProxy, k8sCoreClient core.Interface, metrics *metrics.Metrics, w http.ResponseWriter, r *http.Request, resp *http.Response, trace opentracing.Span) error {
	logger := logging.WithContext(ctx)

	for _, step := range *f {
    logger.With(
  		zap.String("step.name", step.getName()),
  	).Debug("playing step")
		err := step.do(ctx, proxy, k8sCoreClient, metrics, w, r, resp, trace)
		if err == nil {
			continue
		}
		trace.SetTag(tags.Error, true)
		trace.LogKV(
			"event", tags.Error,
			"error.message", err.Error(),
		)
		return err
	}
	return nil
}