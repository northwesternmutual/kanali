package version

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseVersion(t *testing.T) {
	version, commit = "foo", "bar"
	assert.Equal(t, "foo (bar)", parseVersion())
}
