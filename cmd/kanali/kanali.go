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
	"context"
	"fmt"
	"runtime"
	"strings"

	"github.com/northwesternmutual/kanali/cmd/kanali/app"
	"github.com/northwesternmutual/kanali/cmd/kanali/app/options"
	"github.com/northwesternmutual/kanali/pkg/logging"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var version string
var commit string

var rootCmd = &cobra.Command{
	Use:   "kanali",
	Short: "kubernetes native api gateway",
	Long:  "kubernetes native api gateway",
}

var startCmd = &cobra.Command{
	Use:   `start`,
	Short: `start kanali`,
	Long:  `start kanali`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := app.Run(context.Background()); err != nil {
			logging.WithContext(nil).Fatal(err.Error())
		}
	},
}

var versionCmd = &cobra.Command{
	Use:   `version`,
	Short: `version`,
	Long:  `kanali version`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(fmt.Sprintf("%s (%s)", version, commit))
	},
}

func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.kanali")
	viper.AddConfigPath("/etc/kanali/")
	viper.AutomaticEnv()
	viper.SetEnvPrefix("kanali")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.ReadInConfig()

	if err := options.KanaliOptions.AddAll(startCmd); err != nil {
		logging.WithContext(nil).Fatal(err.Error())
	}

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(startCmd)
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	if err := rootCmd.Execute(); err != nil {
		logging.WithContext(nil).Fatal(err.Error())
	}
}
