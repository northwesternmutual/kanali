package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersionCmdInit(t *testing.T) {
	assert.Equal(t, len(RootCmd.Commands()), 2)
	assert.Equal(t, RootCmd.Commands()[1], versionCmd)
}

func TestVersionCmdRun(t *testing.T) {
	org := out
	out = new(bytes.Buffer)
	defer func() { out = org }()

	versionCmdRun(nil, nil)
	assert.Equal(t, out.(*bytes.Buffer).String(), "changeme\n")
}
