package log

import (
	"fmt"

	"go.uber.org/zap"
)

type Level string

var (
	// DebugLevel logs are typically voluminous, and are usually disabled in
	// production.
	DebugLevel Level = Level(zap.DebugLevel.String())
	// InfoLevel is the default logging priority.
	InfoLevel Level = Level(zap.InfoLevel.String())
	// WarnLevel logs are more important than Info, but don't need individual
	// human review.
	WarnLevel Level = Level(zap.WarnLevel.String())
	// ErrorLevel logs are high-priority. If an application is running smoothly,
	// it shouldn't generate any error-level logs.
	ErrorLevel Level = Level(zap.ErrorLevel.String())
	// PanicLevel logs a message, then panics.
	PanicLevel Level = Level(zap.PanicLevel.String())
	// FatalLevel logs a message, then calls os.Exit(1).
	FatalLevel Level = Level(zap.FatalLevel.String())
)

// String implements flag.Value and string.Stringer
func (l *Level) String() string {
	return string(*l)
}

// Set implements flag.Value
func (l *Level) Set(s string) error {
	switch s {
	case "debug", "DEBUG":
		*l = DebugLevel
	case "info", "INFO", "": // make the zero value useful
		*l = InfoLevel
	case "warn", "WARN":
		*l = WarnLevel
	case "error", "ERROR":
		*l = ErrorLevel
	case "panic", "PANIC":
		*l = PanicLevel
	case "fatal", "FATAL":
		*l = FatalLevel
	default:
		return fmt.Errorf("unrecognized level %s", s)
	}
	return nil
}

func (l *Level) Type() string {
	return "github.com/northwesternmutual/kanali/pkg/log.Level"
}
