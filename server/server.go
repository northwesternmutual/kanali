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

package server

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/armon/go-proxyproto"
	"github.com/northwesternmutual/kanali/config"
	h "github.com/northwesternmutual/kanali/handlers"
	"github.com/northwesternmutual/kanali/monitor"
	"github.com/spf13/viper"
)

// Start will start the HTTP server for the Kanali gateway
// It could either be an HTTP or HTTPS server depending on the configuration
func Start(influxCtlr *monitor.InfluxController) {

	var listener net.Listener
	var lerr error
	var scheme string

	router := h.Logger(h.Handler{InfluxController: influxCtlr, H: h.IncomingRequest})

	address := fmt.Sprintf("%s:%d",
		viper.GetString(config.FlagServerBindAddress.GetLong()),
		getKanaliPort(),
	)

	server := &http.Server{Addr: address, Handler: router}

	if viper.GetString(config.FlagTLSCertFile.GetLong()) == "" || viper.GetString(config.FlagTLSKeyFile.GetLong()) == "" {
		scheme = "http"
		listener, lerr = net.Listen("tcp4", address)
		if lerr != nil {
			logrus.Fatalf("error creating http listener: %s", lerr.Error())
			os.Exit(1)
		}
	} else {
		scheme = "https"
		cert, err := tls.LoadX509KeyPair(viper.GetString(config.FlagTLSCertFile.GetLong()), viper.GetString(config.FlagTLSKeyFile.GetLong()))
		if err != nil {
			logrus.Fatal("could not load server cert/key pair")
			os.Exit(1)
		}
		listener, lerr = tls.Listen("tcp4", address, &tls.Config{Certificates: []tls.Certificate{cert}, Rand: rand.Reader})
		if lerr != nil {
			logrus.Fatal("error creating https net listener")
			os.Exit(1)
		}
		// is bi-direction ssl required
		if viper.GetString(config.FlagTLSCaFile.GetLong()) != "" {
			caCert, err := ioutil.ReadFile(viper.GetString(config.FlagTLSCaFile.GetLong()))
			if err != nil {
				logrus.Fatal(err.Error())
				os.Exit(1)
			}
			// load and set client certificate
			caCertPool := x509.NewCertPool()
			caCertPool.AppendCertsFromPEM(caCert)
			tlsConfig := &tls.Config{
				ClientCAs:  caCertPool,
				ClientAuth: tls.RequireAndVerifyClientCert,
			}
			tlsConfig.BuildNameToCertificate()
			server.TLSConfig = tlsConfig
		}
	}

	if viper.GetBool(config.FlagServerProxyProtocol.GetLong()) {
		listener = &proxyproto.Listener{Listener: listener}
	}

	logrus.Infof(fmt.Sprintf("%s server listening on %s", scheme, address))

	if err := server.Serve(listener); err != nil {
		logrus.Fatal(err.Error())
		os.Exit(1)
	}

}

func getKanaliPort() int {
	if viper.GetInt(config.FlagServerPort.GetLong()) > 0 {
		return viper.GetInt(config.FlagServerPort.GetLong())
	}
	if viper.GetString(config.FlagTLSCertFile.GetLong()) == "" || viper.GetString(config.FlagTLSKeyFile.GetLong()) == "" {
		viper.Set(config.FlagServerPort.GetLong(), 80)
		return 80
	}
	viper.Set(config.FlagServerPort.GetLong(), 443)
	return 443
}
