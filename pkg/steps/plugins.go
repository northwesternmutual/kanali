package steps

import (
	"fmt"
	"net/http"
	"path/filepath"
	pluginPkg "plugin"
	"strings"

	"github.com/northwesternmutual/kanali/cmd/kanali/app/options"
	"github.com/northwesternmutual/kanali/pkg/apis/kanali.io/v2"
	kanaliErrors "github.com/northwesternmutual/kanali/pkg/errors"
	"github.com/northwesternmutual/kanali/pkg/plugin"
	"github.com/spf13/viper"
)

const (
	pluginSymbolName = "Plugin"
)

func getPluginFileName(p v2.Plugin) string {
	if strings.Compare(p.Version, "") != 0 {
		return fmt.Sprintf("%s_%s",
			p.Name,
			p.Version,
		)
	}
	return p.Name
}

func getPlugin(pl v2.Plugin) (*plugin.Plugin, error) {
	path, err := getAbsPath(viper.GetString(options.FlagPluginsLocation.GetLong()))
	if err != nil {
		return nil, kanaliErrors.StatusError{Code: http.StatusInternalServerError, Err: fmt.Errorf("file path %s could not be found", viper.GetString(options.FlagPluginsLocation.GetLong()))}
	}

	plug, err := pluginPkg.Open(fmt.Sprintf("%s/%s.so",
		path,
		getPluginFileName(pl),
	))
	if err != nil {
		return nil, kanaliErrors.StatusError{
			Code: http.StatusInternalServerError,
			Err:  fmt.Errorf("could not open plugin %s: %s", pl.Name, err.Error()),
		}
	}

	symPlug, err := plug.Lookup(pluginSymbolName)
	if err != nil {
		return nil, kanaliErrors.StatusError{
			Code: http.StatusInternalServerError,
			Err:  err,
		}
	}

	var p plugin.Plugin
	p, ok := symPlug.(plugin.Plugin)
	if !ok {
		return nil, kanaliErrors.StatusError{
			Code: http.StatusInternalServerError,
			Err:  fmt.Errorf("plugin %s must implement the Plugin interface", pl.Name),
		}
	}

	return &p, nil
}

func getAbsPath(path string) (string, error) {

	p, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	if p[len(p)-1] == '/' {
		if len(p) < 2 {
			return "", nil
		}
		return p[:len(p)-2], nil
	}

	return p, nil

}
