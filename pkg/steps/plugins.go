// Copyright (c) 2017 Northwestern Mutual.
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
