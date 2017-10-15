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

package cmd

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"os"
	"strings"

	"github.com/northwesternmutual/kanali/config"
	"github.com/northwesternmutual/kanali/controller"
	"github.com/northwesternmutual/kanali/crds"
	"github.com/northwesternmutual/kanali/logging"
	"github.com/northwesternmutual/kanali/monitor"
	"github.com/northwesternmutual/kanali/server"
	"github.com/northwesternmutual/kanali/spec"
	"github.com/northwesternmutual/kanali/tracer"
	"github.com/northwesternmutual/kanali/traffic"
	"github.com/opentracing/opentracing-go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.kanali")
	viper.AddConfigPath("/etc/kanali/")
	viper.AutomaticEnv()
	viper.SetEnvPrefix("kanali")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.ReadInConfig()

	if err := config.Flags.AddAll(startCmd); err != nil {
		panic(err)
	}
	RootCmd.AddCommand(startCmd)
	logging.Init(nil)
}

var startCmd = &cobra.Command{
	Use:   `start`,
	Short: `start Kanali`,
	Long:  `start Kanali`,
	Run: func(cmd *cobra.Command, args []string) {

		logger := logging.WithContext(nil)

		ctlr, err := controller.New()
		if err != nil {
			logger.Fatal(err.Error())
			os.Exit(1)
		}

		if err := loadDecryptionKey(viper.GetString(config.FlagPluginsAPIKeyDecriptionKeyFile.GetLong())); err != nil {
			logger.Fatal(err.Error())
			os.Exit(1)
		}

		if err := crds.CreateCRDs(ctlr.APIExtensionsV1beta1Interface); err != nil {
			logger.Fatal(err.Error())
			os.Exit(1)
		}

		go func() {
			if err := ctlr.Watch(context.Background()); err != nil {
				logger.Fatal(err.Error())
			}
		}()

		etcdCtlr, err := traffic.NewController()
		if err != nil {
			logger.Fatal(err.Error())
			os.Exit(1)
		}
		defer func() {
			if err := etcdCtlr.Client.Close(); err != nil {
				logger.Warn(err.Error())
			}
		}()

		go etcdCtlr.MonitorTraffic()

		tracer, closer, err := tracer.Jaeger()
		if err != nil {
			logger.Warn(err.Error())
		} else {
			opentracing.SetGlobalTracer(tracer)
			defer func() {
				if err := closer.Close(); err != nil {
					logger.Warn(err.Error())
				}
			}()
		}

		influxCtlr, err := monitor.NewInfluxdbController()
		if err != nil {
			logger.Warn(err.Error())
		} else {
			defer func() {
				if err := influxCtlr.Client.Close(); err != nil {
					logger.Warn(err.Error())
				}
			}()
		}

		if err := server.Start(influxCtlr); err != nil {
			logger.Fatal(err.Error())
			os.Exit(1)
		}

	},
}

func loadDecryptionKey(location string) error {

	// read in private key
	keyBytes, err := ioutil.ReadFile(location)
	if err != nil {
		return err
	}
	// create a pem block from the private key provided
	block, _ := pem.Decode(keyBytes)
	// parse the pem block into a private key
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return err
	}

	spec.APIKeyDecryptionKey = privateKey

	return nil

}
