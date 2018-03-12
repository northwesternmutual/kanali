package flags

import (
	"github.com/spf13/pflag"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/northwesternmutual/kanali/pkg/log"
)

func LogFlags(fs *pflag.FlagSet) {
	params := []zapcore.Field{}
	fs.VisitAll(func(f *pflag.Flag) {
		params = append(params, zap.String(f.Name, f.Value.String()))
	})
	log.WithContext(nil).With(params...).Info("kanali options")
}
