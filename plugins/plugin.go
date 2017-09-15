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

package plugins

import (
	"context"
	"fmt"
	"net/http"
	pluginPkg "plugin"

	"github.com/northwesternmutual/kanali/config"
	"github.com/northwesternmutual/kanali/controller"
	"github.com/northwesternmutual/kanali/metrics"
	"github.com/northwesternmutual/kanali/spec"
	"github.com/northwesternmutual/kanali/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/spf13/viper"
)

const pluginSymbolName = "Plugin"

// Plugin is an interface that is used for every Plugin used by Kanali.
// If external plugins are developed, they also must conform to this interface.
type Plugin interface {
	OnRequest(ctx context.Context, m *metrics.Metrics, proxy spec.APIProxy, ctlr controller.Controller, req *http.Request, span opentracing.Span) error
	OnResponse(ctx context.Context, m *metrics.Metrics, proxy spec.APIProxy, ctlr controller.Controller, req *http.Request, resp *http.Response, span opentracing.Span) error
}

// GetPlugin will use the Go plugin package and extract
// the plugin
func GetPlugin(plugin spec.Plugin) (*Plugin, error) {
	path, err := utils.GetAbsPath(viper.GetString(config.FlagPluginsLocation.GetLong()))
	if err != nil {
		return nil, utils.StatusError{Code: http.StatusInternalServerError, Err: fmt.Errorf("file path %s could not be found", viper.GetString(config.FlagPluginsLocation.GetLong()))}
	}

	plug, err := pluginPkg.Open(fmt.Sprintf("%s/%s.so",
		path,
		plugin.GetFileName(),
	))
	if err != nil {
		return nil, utils.StatusError{
			Code: http.StatusInternalServerError,
			Err:  fmt.Errorf("could not open plugin %s: %s", plugin.Name, err.Error()),
		}
	}

	symPlug, err := plug.Lookup(pluginSymbolName)
	if err != nil {
		return nil, utils.StatusError{
			Code: http.StatusInternalServerError,
			Err:  err,
		}
	}

	var p Plugin
	p, ok := symPlug.(Plugin)
	if !ok {
		return nil, utils.StatusError{
			Code: http.StatusInternalServerError,
			Err:  fmt.Errorf("plugin %s must implement the Plugin interface", plugin.Name),
		}
	}

	return &p, nil
}
