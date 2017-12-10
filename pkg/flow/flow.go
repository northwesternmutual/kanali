package flow

import (
	"context"
	"net/http"

	"github.com/northwesternmutual/kanali/pkg/apis/kanali.io/v2"
	"github.com/northwesternmutual/kanali/pkg/logging"
	"github.com/northwesternmutual/kanali/pkg/metrics"
	"github.com/northwesternmutual/kanali/pkg/steps"
	"github.com/northwesternmutual/kanali/pkg/tags"
	opentracing "github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"k8s.io/client-go/informers/core"
)

type flow []steps.Step

func New() *flow {
	return &flow{}
}

func (f *flow) Add(steps ...steps.Step) {
	for _, s := range steps {
		*f = append(*f, s)
	}
}

func (f *flow) Play(ctx context.Context, proxy *v2.ApiProxy, k8sCoreClient core.Interface, metrics *metrics.Metrics, w http.ResponseWriter, r *http.Request, resp *http.Response, trace opentracing.Span) error {
	logger := logging.WithContext(ctx)

	for _, step := range *f {
		logger.With(
			zap.String("step.name", step.GetName()),
		).Debug("playing step")
		err := step.Do(ctx, proxy, k8sCoreClient, metrics, w, r, resp, trace)
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
