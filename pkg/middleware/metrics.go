package middleware

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/northwesternmutual/kanali/pkg/logging"
	"github.com/northwesternmutual/kanali/pkg/tags"
	"github.com/northwesternmutual/kanali/pkg/utils"
)

type metricsResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

var (
	ipRegex = regexp.MustCompile("^[0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3}$")
)

func (mrw *metricsResponseWriter) WriteHeader(code int) {
	mrw.statusCode = code
	mrw.ResponseWriter.WriteHeader(code)
}

// Metrics is a middleware that will log and report metrics
// corresponding to the current request.
func Metrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		mrw := &metricsResponseWriter{w, http.StatusOK}

		next.ServeHTTP(mrw, r)

		endTime := time.Now()

		logging.WithContext(r.Context()).Info("request details",
			zap.String(tags.HTTPRequestRemoteAddress, getRemoteAddr(r.RemoteAddr)),
			zap.String(tags.HTTPRequestMethod, r.Method),
			zap.String(tags.HTTPRequestURLPath, utils.ComputeURLPath(r.URL)),
			zap.String(tags.HTTPRequestDuration, fmt.Sprintf("%gms", getReqDuration(startTime, endTime))),
			zap.Int(tags.HTTPResponseStatusCode, mrw.statusCode),
		)
	})
}

// getRemoteAddr return a parsed remote address. There is not defined format for
// http.Request.RemoteAddr which is why this function is required.
func getRemoteAddr(addr string) string {
	if len(addr) < 1 {
		return addr
	}
	if potentialIPAddr := strings.Split(addr, ":")[0]; ipRegex.MatchString(potentialIPAddr) {
		return potentialIPAddr
	}
	return addr
}

func getReqDuration(start, finish time.Time) float64 {
	return (float64(finish.UnixNano()-start.UnixNano()) / float64(time.Millisecond))
}
