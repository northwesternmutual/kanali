package steps

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/northwesternmutual/kanali/pkg/apis/kanali.io/v2"
	"github.com/northwesternmutual/kanali/pkg/logging"
	"github.com/northwesternmutual/kanali/pkg/metrics"
	"github.com/northwesternmutual/kanali/pkg/plugin"
	opentracing "github.com/opentracing/opentracing-go"
	"k8s.io/client-go/informers/core"
)

// PluginsOnResponseStep is factory that defines a step responsible for
// executing the on response lifecycle hook for all the defined plugins
type pluginsOnResponseStep struct{}

func NewPluginsOnResponseStep() Step {
	return pluginsOnResponseStep{}
}

// GetName retruns the name of the PluginsOnResponseStep step
func (step pluginsOnResponseStep) GetName() string {
	return "Plugin OnResponse"
}

// Do executes the logic of the PluginsOnResponseStep step
func (step pluginsOnResponseStep) Do(ctx context.Context, proxy *v2.ApiProxy, k8sCoreClient core.Interface, m *metrics.Metrics, w http.ResponseWriter, r *http.Request, resp *http.Response, trace opentracing.Span) error {

	for _, plugin := range proxy.Spec.Plugins {
		p, err := getPlugin(plugin)
		if err != nil {
			return err
		}
		if err := doOnResponse(ctx, m, plugin, *proxy, r, resp, trace, *p); err != nil {
			return err
		}
	}

	return nil
}

func doOnResponse(ctx context.Context, m *metrics.Metrics, plugin v2.Plugin, proxy v2.ApiProxy, req *http.Request, resp *http.Response, span opentracing.Span, p plugin.Plugin) (e error) {
	logger := logging.WithContext(ctx)

	defer func() {
		if r := recover(); r != nil {
			logger.Error(fmt.Sprintf("OnResponse paniced: %v", r))
			e = errors.New("OnResponse paniced")
		}
	}()

	sp := opentracing.StartSpan(fmt.Sprintf("PLUGIN: ON_RESPONSE: %s", plugin.Name), opentracing.ChildOf(span.Context()))
	defer sp.Finish()

	return p.OnResponse(ctx, plugin.Config, m, proxy, req, resp, sp)
}
