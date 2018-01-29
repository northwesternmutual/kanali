package flags

import (
	"errors"
	"github.com/spf13/viper"
	"strings"

	"github.com/northwesternmutual/kanali/pkg/log"
)

func InitViper(appName string) error {
	if len(appName) < 1 {
		return errors.New("app name undefined")
	}

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/." + appName)
	viper.AddConfigPath("/etc/" + appName + "/")
	viper.AutomaticEnv()
	viper.SetEnvPrefix(appName)
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	if err := viper.ReadInConfig(); err != nil {
		log.WithContext(nil).Warn(err.Error())
	}

	return nil
}
