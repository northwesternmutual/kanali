package steps

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"

	"github.com/northwesternmutual/kanali/cmd/kanali/app/options"
	"github.com/northwesternmutual/kanali/pkg/apis/kanali.io/v2"
	kanaliErrors "github.com/northwesternmutual/kanali/pkg/errors"
	"github.com/northwesternmutual/kanali/pkg/metrics"
	store "github.com/northwesternmutual/kanali/pkg/store/kanali/v2"
	"github.com/northwesternmutual/kanali/pkg/utils"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/spf13/viper"
	"k8s.io/client-go/informers/core"
)

type mockTargetStep struct{}

func NewMockTargetStep() Step {
	return mockTargetStep{}
}

func (step mockTargetStep) GetName() string {
	return "mock target"
}

func (step mockTargetStep) Do(ctx context.Context, proxy *v2.ApiProxy, k8sCoreClient core.Interface, m *metrics.Metrics, w http.ResponseWriter, r *http.Request, resp *http.Response, trace opentracing.Span) error {

	if !mockTargetDefined(proxy) || !mockTargetEnabled(proxy) {
		return nil
	}

	targetPath := utils.ComputeTargetPath(proxy.Spec.Source.Path, proxy.Spec.Target.Path, utils.ComputeURLPath(r.URL))

	mr := store.MockTargetStore().Get(proxy.ObjectMeta.Namespace, proxy.Spec.Target.Backend.Mock.MockTargetName, targetPath, r.Method)
	if mr == nil {
		return &kanaliErrors.StatusError{Code: http.StatusNotFound, Err: errors.New("mock target not found")}
	}

	upstreamHeaders := http.Header{}
	for k, v := range mr.Headers {
		upstreamHeaders.Add(k, v)
	}

	// create a fake response
	responseRecorder := &httptest.ResponseRecorder{
		Code:      mr.StatusCode,
		Body:      bytes.NewBuffer(mr.Body),
		HeaderMap: upstreamHeaders,
	}

	m.Add(metrics.Metric{Name: "http_response_code", Value: strconv.Itoa(mr.StatusCode), Index: true})

	*resp = *responseRecorder.Result()

	return nil

}

func mockTargetDefined(proxy *v2.ApiProxy) bool {
	return len(proxy.Spec.Target.Backend.Mock.MockTargetName) > 0
}

func mockTargetEnabled(proxy *v2.ApiProxy) bool {
	return viper.GetBool(options.FlagProxyEnableMockResponses.GetLong())
}
