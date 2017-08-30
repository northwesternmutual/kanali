package server

import (
	"testing"

	"github.com/northwesternmutual/kanali/config"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestGetKanaliPort(t *testing.T) {
	assert.Equal(t, getKanaliPort(), 80)

	viper.Set(config.FlagKanaliPort.GetLong(), 12345)
	assert.Equal(t, getKanaliPort(), 12345)

	viper.Set(config.FlagKanaliPort.GetLong(), 0)
	viper.Set(config.FlagTLSCertFile.GetLong(), "hi")
	viper.Set(config.FlagTLSPrivateKeyFile.GetLong(), "bye")
	assert.Equal(t, getKanaliPort(), 443)
}
