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

package main

import (
	"math/rand"
	"os"
	"runtime"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	_ "github.com/apache/thrift/lib/go/thrift"
	"github.com/northwesternmutual/kanali/cmd/kanali/app"
	"github.com/northwesternmutual/kanali/cmd/kanali/app/options"
	"github.com/northwesternmutual/kanali/pkg/flags"
	"github.com/northwesternmutual/kanali/pkg/logging"
	"github.com/northwesternmutual/kanali/pkg/server"
	"github.com/northwesternmutual/kanali/pkg/version"
)

const (
	appName             = "kanali"
	appDescriptionShort = "Kubernetes native API gateway."
	appDescriptionLong  = appDescriptionShort
)

var rootCmd = &cobra.Command{
	Use:   appName,
	Short: appDescriptionShort,
	Long:  appDescriptionLong,
}

var startCmd = &cobra.Command{
	Use:   `start`,
	Short: `Start the API gateway`,
	Long:  `Start the API gateway`,
	Run:   startCmdRun,
}

func startCmdRun(cmd *cobra.Command, args []string) {
	ctx := server.SetupSignalHandler()
  logging.Init(nil, viper.GetString(options.FlagProcessLogLevel.GetLong()))
	if err := app.Run(ctx); err != nil {
		logging.WithContext(nil).Error(err.Error())
		os.Exit(1)
	}
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	runtime.GOMAXPROCS(runtime.NumCPU())

	if err := flags.InitViper(appName); err != nil {
		logging.WithContext(nil).Fatal(err.Error())
		os.Exit(1)
	}

	if err := options.KanaliGatewayOptions.AddAll(startCmd); err != nil {
		logging.WithContext(nil).Fatal(err.Error())
		os.Exit(1)
	}

	rootCmd.AddCommand(version.Command())
	rootCmd.AddCommand(startCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
