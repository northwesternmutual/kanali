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
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/northwesternmutual/kanali/cmd/kanali/app/options"
	"github.com/northwesternmutual/kanali/pkg/apis/kanali.io/v2"
	kanaliErrors "github.com/northwesternmutual/kanali/pkg/errors"
	"github.com/northwesternmutual/kanali/pkg/logging"
	"github.com/northwesternmutual/kanali/pkg/metrics"
	tags "github.com/northwesternmutual/kanali/pkg/tags"
	"github.com/northwesternmutual/kanali/pkg/utils"
	opentracing "github.com/opentracing/opentracing-go"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"k8s.io/client-go/informers/core"
)

type httpHandlerFunc func(ctx context.Context, proxy *v2.ApiProxy, k8sCoreClient core.Interface, m *metrics.Metrics, w http.ResponseWriter, r *http.Request, trace opentracing.Span) error

type httpHandler struct {
	*influxController
	k8sCoreClient core.Interface
	httpHandlerFunc
}

func startHTTP(ctx context.Context, handler http.HandlerFunc) error {
	logger := logging.WithContext(nil)

	address := fmt.Sprintf("%s:%d", viper.GetString(options.FlagServerBindAddress.GetLong()), getPort())

	server := &http.Server{Addr: address, Handler: handler}

	if shouldServeInsecurely() {
		logger.Info(fmt.Sprintf("http server listening on %s", address))
		return server.ListenAndServe()
	}

	if !isUsingCustomCACert() {
		logger.Info(fmt.Sprintf("https server listening on %s", address))
		return server.ListenAndServeTLS(viper.GetString(options.FlagTLSCertFile.GetLong()), viper.GetString(options.FlagTLSKeyFile.GetLong()))
	}

	if err := loadCustomCACert(server); err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("https server listening on %s", address))
	return server.ListenAndServeTLS(viper.GetString(options.FlagTLSCertFile.GetLong()), viper.GetString(options.FlagTLSKeyFile.GetLong()))
}

func loadCustomCACert(server *http.Server) error {
	caCert, err := ioutil.ReadFile(viper.GetString(options.FlagTLSCaFile.GetLong()))
	if err != nil {
		return err
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	tlsConfig := &tls.Config{
		ClientCAs:  caCertPool,
		ClientAuth: tls.RequireAndVerifyClientCert,
	}
	tlsConfig.BuildNameToCertificate()
	server.TLSConfig = tlsConfig
	return nil
}

func isUsingCustomCACert() bool {
	return len(viper.GetString(options.FlagTLSCaFile.GetLong())) > 0
}

func shouldServeInsecurely() bool {
	return len(viper.GetString(options.FlagTLSCertFile.GetLong())) < 1 || len(viper.GetString(options.FlagTLSKeyFile.GetLong())) < 1
}

func (h httpHandler) serveHTTP(w http.ResponseWriter, r *http.Request) {
	rqCtx := logging.NewContext(context.Background(), zap.Stringer("correlation_id", uuid.NewV4()))
	logger := logging.WithContext(rqCtx)

	t0 := time.Now()
	m := &metrics.Metrics{}

	defer func() {
		m.Add(
			metrics.Metric{Name: "total_time", Value: int(time.Now().Sub(t0) / time.Millisecond), Index: false},
			metrics.Metric{Name: "http_method", Value: r.Method, Index: true},
			metrics.Metric{Name: "http_uri", Value: utils.ComputeURLPath(r.URL), Index: false},
			metrics.Metric{Name: "client_ip", Value: strings.Split(r.RemoteAddr, ":")[0], Index: false},
		)
		logger.Info("request details",
			zap.String(tags.HTTPRequestRemoteAddress, strings.Split(r.RemoteAddr, ":")[0]),
			zap.String(tags.HTTPRequestMethod, r.Method),
			zap.String(tags.HTTPRequestURLPath, utils.ComputeURLPath(r.URL)),
		)
		go func() {
			if err := h.influxController.writeRequestData(m); err != nil {
				logger.Warn(err.Error())
			} else {
				logger.Debug("wrote metrics to InfluxDB")
			}
		}()
	}()

	sp := opentracing.StartSpan(fmt.Sprintf("%s %s",
		r.Method,
		r.URL.EscapedPath(),
	))
	defer sp.Finish()

	hydrateSpanFromRequest(r, sp)

	err := h.httpHandlerFunc(rqCtx, &v2.ApiProxy{}, h.k8sCoreClient, m, w, r, sp)
	if err == nil {
		return
	}

	var e kanaliErrors.Error
	if _, ok := err.(kanaliErrors.Error); !ok {
		e = kanaliErrors.StatusError{Err: errors.New("unknown error"), Code: http.StatusInternalServerError}
	} else {
		e = err.(kanaliErrors.Error)
	}

	sp.SetTag(tags.HTTPResponseStatusCode, e.Status())
	m.Add(metrics.Metric{Name: "http_response_code", Value: strconv.Itoa(e.Status()), Index: true})
	logger.Info(err.Error(),
		zap.String(tags.HTTPRequestMethod, r.Method),
		zap.String(tags.HTTPRequestURLPath, r.URL.EscapedPath()),
	)

	errStatus, err := json.Marshal(kanaliErrors.JSONErr{Code: e.Status(), Msg: e.Error()})
	if err != nil {
		logger.Warn(err.Error())
	} else {
		sp.SetTag(tags.HTTPResponseBody, string(errStatus))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(e.Status())

	if err := json.NewEncoder(w).Encode(kanaliErrors.JSONErr{Code: e.Status(), Msg: e.Error()}); err != nil {
		logger.Error(err.Error())
	}
}

func getHTTPHandlerFunc(k8sCoreClient core.Interface) httpHandlerFunc {
	return func(ctx context.Context, proxy *v2.ApiProxy, k8sCoreClient core.Interface, m *metrics.Metrics, w http.ResponseWriter, r *http.Request, trace opentracing.Span) error {
		f := &flow{}

		f.add(
			validateProxyStep{},
			pluginsOnRequestStep{},
			mockTargetStep{},
			proxyPassStep{},
			pluginsOnResponseStep{},
			writeResponseStep{},
		)

		return f.play(ctx, proxy, k8sCoreClient, m, w, r, &http.Response{}, trace)
	}
}

func getHTTPHandler(influxCtlr *influxController, k8sCoreClient core.Interface) http.HandlerFunc {
	handler := httpHandler{influxController: influxCtlr, k8sCoreClient: k8sCoreClient, httpHandlerFunc: getHTTPHandlerFunc(k8sCoreClient)}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler.serveHTTP(w, r)
	})
}

func getPort() int {
	if viper.GetInt(options.FlagServerPort.GetLong()) > 0 {
		return viper.GetInt(options.FlagServerPort.GetLong())
	}
	if viper.GetString(options.FlagTLSCertFile.GetLong()) == "" || viper.GetString(options.FlagTLSKeyFile.GetLong()) == "" {
		viper.Set(options.FlagServerPort.GetLong(), 80)
		return 80
	}
	viper.Set(options.FlagServerPort.GetLong(), 443)
	return 443
}
