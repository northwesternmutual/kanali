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
	"context"
	"errors"
	"fmt"
	"net/http"
	pluginPkg "plugin"

	"github.com/northwesternmutual/kanali/controller"
	"github.com/northwesternmutual/kanali/plugins"
	"github.com/northwesternmutual/kanali/spec"
	"github.com/northwesternmutual/kanali/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/spf13/viper"
)

// PluginsOnResponseStep is factory that defines a step responsible for
// executing the on response lifecycle hook for all the defined plugins
type PluginsOnResponseStep struct{}

// GetName retruns the name of the PluginsOnResponseStep step
func (step PluginsOnResponseStep) GetName() string {
	return "Plugin OnResponse"
}

// Do executes the logic of the PluginsOnResponseStep step
func (step PluginsOnResponseStep) Do(ctx context.Context, c *controller.Controller, w http.ResponseWriter, r *http.Request, resp *http.Response, trace opentracing.Span) error {

	// retreive the current proxy which holds
	// the plugins that should be executed, if any
	untypedProxy, err := spec.ProxyStore.Get(r.URL.Path)
	if err != nil || untypedProxy == nil {
		return utils.StatusError{Code: http.StatusNotFound, Err: errors.New("proxy not found")}
	}

	proxy, ok := untypedProxy.(spec.APIProxy)
	if !ok {
		return utils.StatusError{Code: http.StatusNotFound, Err: errors.New("proxy not found")}
	}

	// iterate over all of the plugins associated with this proxy
	for _, plugin := range proxy.Spec.Plugins {

		sp := opentracing.StartSpan(fmt.Sprintf("PLUGIN: ON_RESPONSE: %s", plugin.Name), opentracing.ChildOf(trace.Context()))

		path, err := utils.GetAbsPath(viper.GetString("plugins-location"))
		if err != nil {
			return utils.StatusError{Code: http.StatusInternalServerError, Err: fmt.Errorf("file path %s could not be found", viper.GetString("plugins-path"))}
		}

		// open the plugin module using the standard Go
		// plugin package.
		plug, err := pluginPkg.Open(fmt.Sprintf("%s/%s.so",
			path,
			plugin.GetFileName(),
		))
		if err != nil {
			return utils.StatusError{
				Code: http.StatusInternalServerError,
				Err:  fmt.Errorf("could not open plugin %s", plugin.Name),
			}
		}

		// lookup an exported variable from our plugin module
		symPlug, err := plug.Lookup("Plugin")
		if err != nil {
			return utils.StatusError{
				Code: http.StatusInternalServerError,
				Err:  err,
			}
		}

		var p plugins.Plugin
		p, ok := symPlug.(plugins.Plugin)
		if !ok {
			return utils.StatusError{
				Code: http.StatusInternalServerError,
				Err:  fmt.Errorf("plugin %s must implement the Plugin interface", plugin.Name),
			}
		}

		// execute the OnResponse function of our current plugin
		if err := p.OnResponse(ctx, proxy, *c, r, resp, sp); err != nil {
			return err
		}

		sp.Finish()

	}

	// all plugins, if any, have been executed
	// without any error
	return nil

}
