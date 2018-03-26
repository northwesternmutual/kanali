package pool

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetBuffer(t *testing.T) {
	b := GetBuffer()
	assert.NotNil(t, b)
	b.Write([]byte("foo"))
	assert.Zero(t, GetBuffer().Len())
}

func TestPutBuffer(t *testing.T) {
	b := GetBuffer()
	b.Write([]byte("foo"))
	assert.True(t, b.Len() > 0)
	PutBuffer(b)
	assert.Equal(t, 0, b.Len())
}
