package middleware

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"
  "strconv"

	"go.uber.org/zap"

	"github.com/northwesternmutual/kanali/pkg/log"
	"github.com/northwesternmutual/kanali/pkg/tags"
	"github.com/northwesternmutual/kanali/pkg/utils"
  "github.com/northwesternmutual/kanali/pkg/metrics"
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

    metrics.RequestInFlightCount.Inc()

		next.ServeHTTP(mrw, r)

    metrics.RequestInFlightCount.Dec()

		endTime := time.Now()

		log.WithContext(r.Context()).Info("request details",
			zap.String(tags.HTTPRequestRemoteAddress, getRemoteAddr(r.RemoteAddr)),
			zap.String(tags.HTTPRequestMethod, r.Method),
			zap.String(tags.HTTPRequestURLPath, utils.ComputeURLPath(r.URL)),
			zap.String(tags.HTTPRequestDuration, fmt.Sprintf("%gms", getReqDuration(startTime, endTime))),
			zap.Int(tags.HTTPResponseStatusCode, mrw.statusCode),
		)

    metrics.RequestLatency.WithLabelValues(r.Method).Observe(getReqDuration(startTime, endTime))
    metrics.RequestCount.WithLabelValues(strconv.Itoa(mrw.statusCode), r.Method).Inc()
    if isError(mrw.statusCode) {
      metrics.RequestErrorCount.WithLabelValues(strconv.Itoa(mrw.statusCode), r.Method).Inc()
    }
	})
}

func isError(code int) bool {
  return code > 399 && code < 599
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
