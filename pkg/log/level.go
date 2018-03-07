// Copyright (c) 2018 Northwestern Mutual.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

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

// String implements pflag.Value and string.Stringer
func (l *Level) String() string {
	return string(*l)
}

// Set implements pflag.Value
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

// Type implements pflag.Value
func (l *Level) Type() string {
	return "github.com/northwesternmutual/kanali/pkg/log.Level"
}
