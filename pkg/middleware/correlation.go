package middleware

import (
	"net/http"

	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"

	"github.com/northwesternmutual/kanali/pkg/log"
	"github.com/northwesternmutual/kanali/pkg/tags"
)

// Correlation is a middleware that injects a correlation id into
// the request's context. This middleware is most effective if
// execeted before other middleware.
func Correlation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := log.NewContext(r.Context(), zap.Stringer(tags.HTTPRequestCorrelationId, uuid.NewV4()))
		log.WithContext(ctx).Debug("established new correlation id for this request")
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
