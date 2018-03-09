package tracer

import (
	"fmt"

	"github.com/northwesternmutual/kanali/pkg/log"
	"go.uber.org/zap"
)

type customLogger struct{}

func (l customLogger) Error(msg string) {
	log.WithContext(nil).With(
		zap.String("component", "tracer"),
	).Error(msg)
}

func (l customLogger) Infof(msg string, args ...interface{}) {
	log.WithContext(nil).With(
		zap.String("component", "tracer"),
	).Info(fmt.Sprintf(msg, args...))
}
