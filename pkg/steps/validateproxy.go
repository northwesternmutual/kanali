package steps

import (
	"context"
	"errors"
	"net/http"

	"github.com/northwesternmutual/kanali/pkg/apis/kanali.io/v2"
	kanaliErrors "github.com/northwesternmutual/kanali/pkg/errors"
	"github.com/northwesternmutual/kanali/pkg/metrics"
	store "github.com/northwesternmutual/kanali/pkg/store/kanali/v2"
	"github.com/northwesternmutual/kanali/pkg/tags"
	"github.com/northwesternmutual/kanali/pkg/utils"
	opentracing "github.com/opentracing/opentracing-go"
	"k8s.io/client-go/informers/core"
)

// ValidateProxyStep is factory that defines a step responsible for
// validating that an incoming request matches a proxy that Kanali
// has stored in memory
type validateProxyStep struct{}

func NewValidateProxyStep() Step {
	return validateProxyStep{}
}

// GetName retruns the name of the ValidateProxyStep step
func (step validateProxyStep) GetName() string {
	return "Validate Proxy"
}

// Do executes the logic of the ValidateProxyStep step
func (step validateProxyStep) Do(ctx context.Context, proxy *v2.ApiProxy, k8sCoreClient core.Interface, m *metrics.Metrics, w http.ResponseWriter, r *http.Request, resp *http.Response, trace opentracing.Span) error {

	typedProxy := store.ApiProxyStore().Get(utils.ComputeURLPath(r.URL))
	if typedProxy == nil {
		trace.SetTag(tags.KanaliProxyName, "unknown")
		trace.SetTag(tags.KanaliProxyNamespace, "unknown")

		m.Add(
			metrics.Metric{Name: "proxy_name", Value: "unknown", Index: true},
			metrics.Metric{Name: "proxy_namespace", Value: "unknown", Index: true},
		)

		return kanaliErrors.StatusError{Code: http.StatusNotFound, Err: errors.New("proxy not found")}
	}

	*proxy = *typedProxy

	trace.SetTag(tags.KanaliProxyName, proxy.ObjectMeta.Name)
	trace.SetTag(tags.KanaliProxyNamespace, proxy.ObjectMeta.Namespace)

	m.Add(
		metrics.Metric{Name: "proxy_name", Value: proxy.ObjectMeta.Name, Index: true},
		metrics.Metric{Name: "proxy_namespace", Value: proxy.ObjectMeta.Namespace, Index: true},
	)

	return nil

}
