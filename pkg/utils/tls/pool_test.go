package tls

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSystemCertPool(t *testing.T) {
	pool, err := GetSystemCertPool()
	assert.Nil(t, err)
	assert.NotNil(t, pool)
}
