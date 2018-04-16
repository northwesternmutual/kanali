// Copyright (c) 2018 Northwestern Mutual.
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

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"

	"github.com/northwesternmutual/kanali/pkg/chain"
	"github.com/northwesternmutual/kanali/pkg/client/clientset/versioned"
	"github.com/northwesternmutual/kanali/pkg/flags"
	"github.com/northwesternmutual/kanali/pkg/log"
	_ "github.com/northwesternmutual/kanali/pkg/metrics"
	"github.com/northwesternmutual/kanali/pkg/middleware"
	"github.com/northwesternmutual/kanali/pkg/run"
	"github.com/northwesternmutual/kanali/pkg/server"
	"github.com/northwesternmutual/kanali/pkg/utils"
)

func Run(sigCtx context.Context) error {
	logger := log.WithContext(sigCtx)

	kanaliClientset, err := createClientsets()
	if err != nil {
		logger.Fatal(err.Error())
		return err
	}

	gatewayServer := server.Prepare(&server.Options{
		Name:         "gateway",
		InsecureAddr: viper.GetString(flags.FlagServerInsecureBindAddress.GetLong()),
		SecureAddr:   viper.GetString(flags.FlagServerSecureBindAddress.GetLong()),
		InsecurePort: viper.GetInt(flags.FlagServerInsecurePort.GetLong()),
		SecurePort:   viper.GetInt(flags.FlagServerSecurePort.GetLong()),
		TLSKey:       viper.GetString(flags.FlagServerTLSKeyFile.GetLong()),
		TLSCert:      viper.GetString(flags.FlagServerTLSCertFile.GetLong()),
		TLSCa:        viper.GetString(flags.FlagServerTLSCaFile.GetLong()),
		Handler: chain.New().Add(
			middleware.Correlation,
			middleware.Recover,
			middleware.Metrics,
		).Link(middleware.Validator(kanaliClientset)),
		Logger: logger.Sugar(),
	})

	profilingServer := server.Prepare(&server.Options{
		Name:         "profiling",
		InsecureAddr: viper.GetString(flags.FlagProfilingInsecureBindAddress.GetLong()),
		InsecurePort: viper.GetInt(flags.FlagProfilingInsecurePort.GetLong()),
		Handler:      server.ProfilingHandler(),
		Logger:       logger.Sugar(),
	})

	metricsServer := server.Prepare(&server.Options{
		Name:         "prometheus",
		InsecureAddr: viper.GetString(flags.FlagPrometheusServerInsecureBindAddress.GetLong()),
		InsecurePort: viper.GetInt(flags.FlagPrometheusServerInsecurePort.GetLong()),
		Handler:      promhttp.Handler(),
		Logger:       logger.Sugar(),
	})

	ctx, cancel := context.WithCancel(sigCtx)

	var g run.Group
	g.Add(ctx, metricsServer.IsEnabled(), metricsServer.Name(), metricsServer)
	g.Add(ctx, run.Always, gatewayServer.Name(), gatewayServer)
	g.Add(ctx, profilingServer.IsEnabled(), profilingServer.Name(), profilingServer)
	g.Add(ctx, run.Always, "parent process", run.MonitorContext(cancel))
	return g.Run()
}

func createClientsets() (versioned.Interface, error) {
	config, err := utils.GetKubernetesRestConfig(viper.GetString(flags.FlagKubernetesKubeConfig.GetLong()))
	if err != nil {
		return nil, err
	}

	kanaliClientset, err := versioned.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return kanaliClientset, nil
}
