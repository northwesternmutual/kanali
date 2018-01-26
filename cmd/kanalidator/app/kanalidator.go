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

package app

import (
	"context"

	"github.com/oklog/run"
	"github.com/spf13/viper"

	"github.com/northwesternmutual/kanali/cmd/kanalidator/app/options"
	"github.com/northwesternmutual/kanali/pkg/client/clientset/versioned"
	"github.com/northwesternmutual/kanali/pkg/kanalidator/server"
	"github.com/northwesternmutual/kanali/pkg/logging"
	"github.com/northwesternmutual/kanali/pkg/utils"
)

func Run(sigCtx context.Context) error {
	ctx, cancel := context.WithCancel(sigCtx)
	logger := logging.WithContext(nil)

	config, err := utils.GetRestConfig(viper.GetString(options.FlagKubernetesKubeConfig.GetLong()))
	if err != nil {
		return err
	}

	kanaliClientset, err = versioned.NewForConfig(config)
	if err != nil {
		return err
	}

	kanalidator, err := server.New(kanaliClientset)
	if err != nil {
		logger.Fatal(err.Error())
		return err
	}

	var g run.Group

	g.Add(func() error {
		err := kanalidator.Run()
		logger.Error(err.Error())
		return err
	}, func(error) {
		kanalidator.Close()
	})

	return g.Run()
}
