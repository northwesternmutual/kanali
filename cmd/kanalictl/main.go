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

package main

import (
	"math/rand"
	"os"
	"runtime"
	"time"

	"github.com/spf13/cobra"

	"github.com/northwesternmutual/kanali/cmd/kanalictl/app"
	"github.com/northwesternmutual/kanali/cmd/kanalictl/app/options"
	"github.com/northwesternmutual/kanali/pkg/flags"
	"github.com/northwesternmutual/kanali/pkg/log"
	"github.com/northwesternmutual/kanali/pkg/version"
)

const (
	appName             = "kanalictl"
	appDescriptionShort = "Command line interface for Kanali."
	appDescriptionLong  = appDescriptionShort
)

var rootCmd = &cobra.Command{
	Use:   appName,
	Short: appDescriptionShort,
	Long:  appDescriptionLong,
}

var apiKeyCmd = &cobra.Command{
	Use:   `apikey`,
	Short: `Preforms operations on API key resources`,
	Long:  `Preforms operations on API key resources`,
}

var decryptCmd = &cobra.Command{
	Use:   `decrypt`,
	Short: `Decrypts API key resources`,
	Long:  `Decrypts API key resources`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := app.Decrypt(os.Stdout, os.Stderr); err != nil {
			os.Exit(1)
		}
		os.Exit(0)
	},
}

var generateCmd = &cobra.Command{
	Use:   `generate`,
	Short: `Creates an API key`,
	Long:  `Creates an API key`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := app.Generate(os.Stdout, os.Stderr); err != nil {
			os.Exit(1)
		}
		os.Exit(0)
	},
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	runtime.GOMAXPROCS(runtime.NumCPU())

	if err := flags.NewFlagSet().Add(
		options.FlagRSAPrivateKeyFile,
		options.FlagKeyInFile,
	).AddAll(decryptCmd); err != nil {
		log.WithContext(nil).Fatal(err.Error())
		os.Exit(1)
	}

	if err := flags.NewFlagSet().Add(
		options.FlagRSAPublicKeyFile,
		options.FlagKeyName,
		options.FlagKeyData,
		options.FlagKeyOutFile,
		options.FlagKeyLength,
	).AddAll(generateCmd); err != nil {
		log.WithContext(nil).Fatal(err.Error())
		os.Exit(1)
	}

	apiKeyCmd.AddCommand(decryptCmd)
	apiKeyCmd.AddCommand(generateCmd)

	rootCmd.AddCommand(version.Command())
	rootCmd.AddCommand(apiKeyCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
