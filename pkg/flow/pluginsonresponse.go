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
	"net/http"

	"github.com/northwesternmutual/kanali/pkg/errors"
	"github.com/northwesternmutual/kanali/pkg/logging"
	"github.com/northwesternmutual/kanali/pkg/plugin"
	store "github.com/northwesternmutual/kanali/pkg/store/kanali/v2"
	"github.com/northwesternmutual/kanali/pkg/utils"
)

// PluginsOnResponseStep is factory that defines a step responsible for
// executing the on response lifecycle hook for all the defined plugins
type pluginsOnResponseStep struct{}

func PluginsOnResponseStep() Step {
	return pluginsOnResponseStep{}
}

// GetName retruns the name of the PluginsOnResponseStep step
func (step pluginsOnResponseStep) Name() string {
	return "Plugin OnResponse"
}

// Do executes the logic of the PluginsOnResponseStep step
func (step pluginsOnResponseStep) Do(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	logger := logging.WithContext(r.Context())

	proxy := store.ApiProxyStore().Get(utils.ComputeURLPath(r.URL))
	if proxy == nil {
		logger.Warn(errors.ErrorProxyNotFound.Message)
		return errors.ErrorProxyNotFound
	}

	for _, plugin := range proxy.Spec.Plugins {
		p, err := getPlugin(ctx, plugin)
		if err != nil {
			logger.Error(err.Error())
			return err
		}
		if err := doOnResponse(ctx, p, plugin.Config, w, r); err != nil {
			logger.Error(err.Error())
			return err
		}
	}

	return next()
}

func doOnResponse(ctx context.Context, p plugin.Plugin, config map[string]string, w http.ResponseWriter, r *http.Request) (err error) {
	defer func() {
		if r := recover(); r != nil {
			logging.WithContext(ctx).Error(fmt.Sprintf("%v", r))
			err = errors.ErrorPluginRuntimeError
		}
	}()
	return p.OnResponse(ctx, config, w, r)
}
