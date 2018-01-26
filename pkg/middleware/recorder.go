package middleware

import (
	"net/http"
	"net/http/httptest"

	"github.com/northwesternmutual/kanali/pkg/utils"
)

// Recorder will inject a recorder in place of the actual http response so
// that it can record what all other middleware functions do when mutating
// the http response. This will prevent any other middleware from writing
// to the actually http response.
func Recorder(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec := httptest.NewRecorder()
		next.ServeHTTP(rec, r)
		utils.TransferResponse(rec, w)
	})
}
