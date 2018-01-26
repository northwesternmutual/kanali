package middleware

import (
	"context"
	"net/http"

	"github.com/northwesternmutual/kanali/pkg/errors"
	"github.com/northwesternmutual/kanali/pkg/flow"
	"github.com/northwesternmutual/kanali/pkg/logging"
	"github.com/northwesternmutual/kanali/pkg/tags"
	"github.com/northwesternmutual/kanali/pkg/tracer"
	opentracing "github.com/opentracing/opentracing-go"
)

// Gateway is an http.HandlerFunc that preforms the core functions of the Kanali gateway.
// This is meant to be the final http.Handler if multiple middlewares are to be used.
func Gateway(w http.ResponseWriter, r *http.Request) {
	logger := logging.WithContext(r.Context())
	span := tracer.StartSpan(r)
	ctx := opentracing.ContextWithSpan(context.Background(), span)
	defer span.Finish()

	tracer.HydrateSpanFromRequest(r, span)

	if err := flow.New().Add(
		flow.ValidateProxyStep(),
		flow.PluginsOnRequestStep(),
		flow.MockTargetStep(),
		flow.ProxyPassStep(),
		flow.PluginsOnResponseStep(),
	).Play(ctx, w, r); err != nil {
		err, data := errors.ToJSON(err)
		span.SetTag(tags.HTTPResponseBody, string(data))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(err.Status)
		if _, err := w.Write(data); err != nil {
			logger.Error(err.Error())
		}
	}

	// TODO:
	//tracer.HydrateSpanFromResponse(r, span)
}
