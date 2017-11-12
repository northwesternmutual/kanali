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
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"os"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/northwesternmutual/kanali/config"
	"github.com/northwesternmutual/kanali/controller"
	"github.com/northwesternmutual/kanali/monitor"
	"github.com/northwesternmutual/kanali/server"
	"github.com/northwesternmutual/kanali/spec"
	"github.com/northwesternmutual/kanali/tracer"
	"github.com/opentracing/opentracing-go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {

	if err := config.Flags.AddAll(startCmd); err != nil {
		logrus.Fatalf("could not add flag to command: %s", err.Error())
		os.Exit(1)
	}

	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.kanali")
	viper.AddConfigPath("/etc/kanali/")

	viper.AutomaticEnv()
	viper.SetEnvPrefix("kanali")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	if err := viper.ReadInConfig(); err != nil {
		logrus.Warn("couldn't find any config file, using env variables and/or cli flags")
	}

	RootCmd.AddCommand(startCmd)
}

var startCmd = &cobra.Command{
	Use:   `start`,
	Short: `start Kanali`,
	Long:  `start Kanali`,
	Run: func(cmd *cobra.Command, args []string) {

		// set logging level
		if level, err := logrus.ParseLevel(viper.GetString(config.FlagProcessLogLevel.GetLong())); err != nil {
			logrus.SetLevel(logrus.InfoLevel)
			logrus.Info("could not parse logging level")
		} else {
			logrus.SetLevel(level)
		}

		// create new k8s controller
		ctlr, err := controller.New()
		if err != nil {
			logrus.Fatalf("could not create controller: %s", err.Error())
			os.Exit(1)
		}

		// load decryption key into memory
		if err := loadDecryptionKey(viper.GetString(config.FlagPluginsAPIKeyDecriptionKeyFile.GetLong())); err != nil {
			logrus.Fatalf("could not load decryption key: %s", err.Error())
			os.Exit(1)
		}

		// create tprs
		if err := ctlr.CreateTPRs(); err != nil {
			logrus.Fatalf("could not create TPRs: %s", err.Error())
			os.Exit(1)
		}

		go ctlr.Watch()

		// start UDP server
		go func() {
			if err := server.StartUDPServer(); err != nil {
				logrus.Fatal(err.Error())
				os.Exit(1)
			}
		}()

		tracer, closer, err := tracer.Jaeger()
		if err != nil {
			logrus.Warnf("error create Jaeger tracer: %s", err.Error())
		} else {
			opentracing.SetGlobalTracer(tracer)
			defer func() {
				if err := closer.Close(); err != nil {
					logrus.Warnf("error closing Jaeger tracer: %s", err.Error())
				}
			}()
		}

		influxCtlr, err := monitor.NewInfluxdbController()
		if err != nil {
			logrus.Warnf("error connecting to InfluxDB: %s", err.Error())
		} else {
			go influxCtlr.Run()
			defer func() {
				if err := influxCtlr.Client.Close(); err != nil {
					logrus.Warnf("error closing the connection to InfluxDB: %s", err.Error())
				}
			}()
		}

		server.Start(influxCtlr)

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
