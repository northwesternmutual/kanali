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

package flow

import (
	"context"
	"fmt"
	"path/filepath"
	pluginPkg "plugin"

	"github.com/northwesternmutual/kanali/cmd/kanali/app/options"
	"github.com/northwesternmutual/kanali/pkg/apis/kanali.io/v2"
	"github.com/northwesternmutual/kanali/pkg/errors"
	"github.com/northwesternmutual/kanali/pkg/logging"
	"github.com/northwesternmutual/kanali/pkg/plugin"
	"github.com/spf13/viper"
)

const (
	// Lookup searches for a symbol named symName in plugin p. A symbol is any exported variable or function.
	// It reports an error if the symbol is not found. It is safe for concurrent use by multiple goroutines.
	symName = "Plugin"
)

// join will return the path to the compiled plugin.
// While this path is not guaranteed to be absolute,
// the plugin package will make it absolute:
// https://github.com/golang/go/blob/master/src/plugin/plugin_dlopen.go#L56-L60
func combinePath(basePath string, plugin v2.Plugin) string {
	name := plugin.Name
	if len(plugin.Version) > 0 {
		name += "_" + plugin.Version
	}
	return filepath.Join(basePath, name+".so")
}

func getPlugin(ctx context.Context, pl v2.Plugin) (plugin.Plugin, error) {
	basePath := viper.GetString(options.FlagPluginsLocation.GetLong())
	logger := logging.WithContext(ctx)

	plug, err := pluginPkg.Open(combinePath(basePath, pl))
	if err != nil {
		logger.Error(err.Error())
		return nil, errors.ErrorCouldNotLoadPlugin
	}

	symbol, err := plug.Lookup(symName)
	if err != nil {
		logger.Error(err.Error())
		return nil, errors.ErrorCouldNotLookupPluginSymbol
	}

	p, ok := symbol.(plugin.Plugin)
	if !ok {
		logger.Error(fmt.Sprintf("plugin %s must implement the Plugin interface", pl.Name))
		return nil, errors.ErrorPluginIncorrectInterface
	}

	return p, nil
}
