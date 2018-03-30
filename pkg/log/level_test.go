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
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	debugLvl = DebugLevel
	infoLvl  = InfoLevel
	warnLvl  = WarnLevel
	errorLvl = ErrorLevel
	panicLvl = PanicLevel
	fatalLvl = FatalLevel
)

func TestString(t *testing.T) {
	assert.Equal(t, "debug", (&debugLvl).String())
	assert.Equal(t, "info", (&infoLvl).String())
	assert.Equal(t, "warn", (&warnLvl).String())
	assert.Equal(t, "error", (&errorLvl).String())
	assert.Equal(t, "panic", (&panicLvl).String())
	assert.Equal(t, "fatal", (&fatalLvl).String())
}

func TestType(t *testing.T) {
	assert.Equal(t, "string", (&debugLvl).Type())
}

func TestSet(t *testing.T) {
	l := new(Level)
	assert.Nil(t, l.Set("debug"))
	assert.Equal(t, DebugLevel, *l)
	assert.Nil(t, l.Set("INFO"))
	assert.Equal(t, InfoLevel, *l)
	assert.Nil(t, l.Set("warn"))
	assert.Equal(t, WarnLevel, *l)
	assert.Nil(t, l.Set("ERROR"))
	assert.Equal(t, ErrorLevel, *l)
	assert.Nil(t, l.Set("panic"))
	assert.Equal(t, PanicLevel, *l)
	assert.Nil(t, l.Set("FATAL"))
	assert.Equal(t, FatalLevel, *l)
	assert.NotNil(t, l.Set("foo"))
}
