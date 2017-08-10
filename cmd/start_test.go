package cmd

import (
	"testing"

	"github.com/northwesternmutual/kanali/config"
	"github.com/stretchr/testify/assert"
)

func TestStartCmdInit(t *testing.T) {
	assert.Equal(t, len(RootCmd.Commands()), 2)
	assert.Equal(t, RootCmd.Commands()[0], startCmd)

	for _, f := range config.Flags {
		assert.Equal(t, startCmd.Flag(f.GetLong()).Name, f.GetLong())
		assert.Equal(t, startCmd.Flag(f.GetLong()).Shorthand, f.GetShort())
		assert.Equal(t, startCmd.Flag(f.GetLong()).Usage, f.GetUsage())
	}
}
