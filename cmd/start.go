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
	"os"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/northwesternmutual/kanali/config"
	"github.com/northwesternmutual/kanali/controller"
	"github.com/northwesternmutual/kanali/monitor"
	"github.com/northwesternmutual/kanali/server"
	"github.com/northwesternmutual/kanali/tracer"
	"github.com/northwesternmutual/kanali/utils"
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
		if level, err := logrus.ParseLevel(viper.GetString("log-level")); err != nil {
			logrus.SetLevel(logrus.InfoLevel)
			logrus.Info("could not parse logging level")
		} else {
			logrus.SetLevel(level)
		}

		// create new k8s controller
		ctlr, err := controller.New()
		if err != nil {
			logrus.Panic(err.Error())
			os.Exit(1)
		}

		// load decryption key into memory
		if err := utils.LoadDecryptionKey(viper.GetString("decryption-key-file")); err != nil {
			logrus.Panic(err.Error())
			os.Exit(1)
		}

		// create tprs
		if err := ctlr.CreateTPRs(); err != nil {
			logrus.Panic(err.Error())
			os.Exit(1)
		}

		// start walking k8s resources
		go func() {
			if err := ctlr.Watch(); err != nil {
				logrus.Fatal(err.Error())
				os.Exit(1)
			}
		}()

		// start UDP server
		go func() {
			if err := server.StartUDPServer(); err != nil {
				logrus.Fatal(err.Error())
				os.Exit(1)
			}
		}()

		// potentially start tracing server
		if viper.GetBool("enable-tracing") {
			tracer, closer, err := tracer.Jaeger()
			if err != nil {
				logrus.Fatal(err.Error())
				os.Exit(1)
			}
			logrus.Infof("starting global tracer")
			opentracing.SetGlobalTracer(tracer)
			defer func() {
				if err := closer.Close(); err != nil {
					logrus.Warnf("there was a problem closing the tracer: %s", err.Error())
				}
			}()
		}

		// attempt to create influxdb client
		influxCtlr, err := monitor.NewInfluxdbController()
		if err != nil {
			logrus.Warnf("there was an error connecting to influxdb - analytics and monitoring will not be available", err.Error())
		} else {
			defer func() {
				if err := influxCtlr.Client.Close(); err != nil {
					logrus.Warnf("there was a problem closing the connection to influxdb: %s", err.Error())
				}
			}()
		}

		// start kanali readiness server
		go server.StatusServer(ctlr)

		// start kanali gateway server
		server.Start(ctlr, influxCtlr)

	},
}
