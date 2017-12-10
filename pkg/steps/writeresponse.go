package steps

import (
	"context"
	"io"
	"net/http"
	"strconv"

	"github.com/northwesternmutual/kanali/pkg/apis/kanali.io/v2"
	"github.com/northwesternmutual/kanali/pkg/metrics"
	"github.com/northwesternmutual/kanali/pkg/tracer"
	opentracing "github.com/opentracing/opentracing-go"
	"k8s.io/client-go/informers/core"
)

// WriteResponseStep is factory that defines a step responsible for writing
// an HTTP response
type writeResponseStep struct{}

func NewWriteResponseStep() Step {
	return writeResponseStep{}
}

// GetName retruns the name of the WriteResponseStep step
func (step writeResponseStep) GetName() string {
	return "Write Response"
}

// Do executes the logic of the WriteResponseStep step
func (step writeResponseStep) Do(ctx context.Context, proxy *v2.ApiProxy, k8sCoreClient core.Interface, m *metrics.Metrics, w http.ResponseWriter, r *http.Request, resp *http.Response, span opentracing.Span) error {

	for k, v := range resp.Header {
		for _, value := range v {
			w.Header().Set(k, value)
		}
	}

	m.Add(metrics.Metric{Name: "http_response_code", Value: strconv.Itoa(resp.StatusCode), Index: true})

	tracer.HydrateSpanFromResponse(resp, span)

	w.WriteHeader(resp.StatusCode)

	if _, err := io.Copy(w, resp.Body); err != nil {
		return err
	}

	return nil
}
