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

type pluginsOnRequestStep struct{}

func NewPluginsOnRequestStep() Step {
	return pluginsOnRequestStep{}
}

func (step pluginsOnRequestStep) GetName() string {
	return "plugin onrequest"
}

func (step pluginsOnRequestStep) Do(ctx context.Context, proxy *v2.ApiProxy, k8sCoreClient core.Interface, m *metrics.Metrics, w http.ResponseWriter, r *http.Request, resp *http.Response, trace opentracing.Span) error {

	for _, plugin := range proxy.Spec.Plugins {
		p, err := getPlugin(plugin)
		if err != nil {
			return err
		}
		if err := doOnRequest(ctx, m, plugin, *proxy, r, trace, *p); err != nil {
			return err
		}
	}
	return nil
}

func doOnRequest(ctx context.Context, m *metrics.Metrics, plugin v2.Plugin, proxy v2.ApiProxy, req *http.Request, span opentracing.Span, p plugin.Plugin) (e error) {
	logger := logging.WithContext(ctx)

	defer func() {
		if r := recover(); r != nil {
			logger.Error(fmt.Sprintf("OnRequest paniced: %v", r))
			e = errors.New("OnRequest paniced")
		}
	}()

	sp := opentracing.StartSpan(fmt.Sprintf("PLUGIN: ON_REQUEST: %s", plugin.Name), opentracing.ChildOf(span.Context()))
	defer sp.Finish()

	return p.OnRequest(ctx, plugin.Config, m, proxy, req, sp)
}
