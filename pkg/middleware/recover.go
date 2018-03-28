package middleware

import (
	"net/http"

	"k8s.io/apimachinery/pkg/util/runtime"

	"github.com/northwesternmutual/kanali/pkg/errors"
	"github.com/northwesternmutual/kanali/pkg/log"
)

func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := log.WithContext(r.Context())

		defer runtime.HandleCrash(func(err interface{}) {
			logger.Error("kanali gateway paniced")

			jsonErr, data := errors.ToJSON(nil)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(jsonErr.Status)
			if _, err := w.Write(data); err != nil {
				logger.Error(err.Error())
			}

		})

		next.ServeHTTP(w, r)
	})
}
