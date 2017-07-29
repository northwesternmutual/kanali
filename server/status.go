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
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/northwesternmutual/kanali/config"
	"github.com/northwesternmutual/kanali/controller"
	"github.com/northwesternmutual/kanali/spec"
	"github.com/spf13/viper"
	"k8s.io/kubernetes/pkg/api"
)

func init() {

	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)

	if level, err := logrus.ParseLevel(viper.GetString(config.FlagLogLevel.GetLong())); err != nil {
		logrus.SetLevel(logrus.InfoLevel)
	} else {
		logrus.SetLevel(level)
	}

}

// StatusServer will start the server used to assess the health of Kanali
func StatusServer(c *controller.Controller) {

	if err := http.ListenAndServe(fmt.Sprintf(":%d", viper.GetInt(config.FlagStatusPort.GetLong())), statusHandler(c)); err != nil {
		logrus.Fatal(err)
		os.Exit(1)
	}

}

func statusHandler(c *controller.Controller) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		switch strings.ToLower(strings.Split(r.URL.Path[1:], "/")[0]) {
		case "liveness":
			if checkLiveness(c) {
				w.WriteHeader(http.StatusOK)
			} else {
				w.WriteHeader(http.StatusBadRequest)
			}
		case "readiness":
			if checkReadiness(c) {
				w.WriteHeader(http.StatusOK)
			} else {
				w.WriteHeader(http.StatusBadRequest)
			}
		default:
			w.WriteHeader(http.StatusBadRequest)
		}

	})

}

// technically, Kanali doesn't need a live connection to the
// k8s apiserver for it to stay alive. It just needs to have
// some proxies loaded in memory. Now there still might be
// lots of reasons why a certain request may not succeed
func checkLiveness(c *controller.Controller) bool {

	if !spec.ProxyStore.IsEmpty() {
		return true
	}

	// proxy store is empty so we'll need to be able to connect to the k8s apiserver
	if _, err := c.ClientSet.Core().Endpoints(api.NamespaceAll).List(api.ListOptions{}); err != nil {
		logrus.Errorf("liveness probe error: %s", err.Error())
		return false
	}

	return true

}

// we don't want a new Kanali pod registered as ready
// if we can't connect to the k8s apiserver so let's check that
func checkReadiness(c *controller.Controller) bool {

	// there's no reason we single out the endpoints api.
	// it happens to be one that Kanali already has permission to access.
	if _, err := c.ClientSet.Core().Endpoints(api.NamespaceAll).List(api.ListOptions{}); err != nil {
		logrus.Errorf("readiness probe error: %s", err.Error())
		return false
	}

	return true

}
