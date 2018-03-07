package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	assert.Equal(t, "debug", (&DebugLevel).String())
	assert.Equal(t, "panic", (&PanicLevel).String())
}

func TestType(t *testing.T) {
	assert.Equal(t, "github.com/northwesternmutual/kanali/pkg/log.Level", (&DebugLevel).Type())
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
