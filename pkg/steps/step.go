package steps

import (
	"context"
	"net/http"

	"github.com/northwesternmutual/kanali/pkg/apis/kanali.io/v2"
	"github.com/northwesternmutual/kanali/pkg/metrics"
	opentracing "github.com/opentracing/opentracing-go"
	"k8s.io/client-go/informers/core"
)

type Step interface {
	GetName() string
	Do(ctx context.Context, proxy *v2.ApiProxy, k8sCoreClient core.Interface, metrics *metrics.Metrics, w http.ResponseWriter, r *http.Request, resp *http.Response, trace opentracing.Span) error
}
