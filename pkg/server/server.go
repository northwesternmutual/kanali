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

package server

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"golang.org/x/sync/errgroup"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
)

type serverParams struct {
	secureServer   *http.Server
	insecureServer *http.Server
	err            []error
	options        *Options
}

// Options contains configuration for http(s) server(s)
type Options struct {
	Name                   string
	InsecureAddr           string
	SecureAddr             string
	InsecurePort           int
	SecurePort             int
	TLSKey, TLSCert, TLSCa string
	Handler                http.Handler
	Logger                 logger
}

type logger interface {
	Info(...interface{})
	Error(...interface{})
}

// Prepare will construct net.Listerer instantiations for the requested
// server(s). If an error is encounted, it will be hidden within the returned
// type for evaluation by that type's methods.
func Prepare(opts *Options) *serverParams {
	f, err := os.Open(opts.TLSCa)
	if err != nil && len(opts.TLSCa) > 0 {
		return &serverParams{err: []error{err}}
	} else {
		defer f.Close()
	}
	return prepare(opts, f)
}

func prepare(opts *Options, r io.Reader) *serverParams {
	params := &serverParams{options: opts}
	insecureAddr, secureAddr := fmt.Sprintf("%s:%d", opts.InsecureAddr, opts.InsecurePort), fmt.Sprintf("%s:%d", opts.SecureAddr, opts.SecurePort)

	if opts.SecurePort > 0 && len(opts.TLSCa) > 0 { // client has requested an HTTPS server with mutual TLS
		tlsConfig, err := getTLSConfigFromReader(r)
		if err != nil {
			params.err = append(params.err, err)
		} else {
			params.secureServer = &http.Server{Addr: secureAddr, TLSConfig: tlsConfig, Handler: opts.Handler}
		}
	} else if opts.SecurePort > 0 && len(opts.TLSCa) < 1 { // client has requested an HTTPS server with one-way TLS
		params.secureServer = &http.Server{Addr: secureAddr, Handler: opts.Handler}
	}
	if opts.InsecurePort > 0 { // client has requested an HTTP server
		params.insecureServer = &http.Server{Addr: insecureAddr, Handler: opts.Handler}
	}
	return params
}

// Run will start http(s) server(s) according to the reciever's configuration.
// If an error occurred durring the receiver's construction, that error will be
// returned immedietaly. Otherwise, run will return the non-nill error from the
// the first server that terminates.
func (params *serverParams) Run(context.Context) error {
	if params.err != nil && len(params.err) > 0 {
		return utilerrors.NewAggregate(params.err)
	}

	var g errgroup.Group
	if params.secureServer != nil {
		params.options.Logger.Info(fmt.Sprintf("starting %s %s server on %s:%d", params.options.Name, "HTTPS", params.options.SecureAddr, params.options.SecurePort))
		g.Go(func() error {
			return params.secureServer.ListenAndServeTLS(params.options.TLSCert, params.options.TLSKey)
		})
	}
	if params.insecureServer != nil {
		params.options.Logger.Info(fmt.Sprintf("starting %s %s server on %s:%d", params.options.Name, "HTTP", params.options.InsecureAddr, params.options.InsecurePort))
		g.Go(func() error {
			return params.insecureServer.ListenAndServe()
		})
	}
	err := g.Wait()
	if err != http.ErrServerClosed {
		params.options.Logger.Error(err.Error())
		return err
	}
	return nil
}

// Close will gracefully terminate http(s) server(s) that were bootstrapped
// according to the reciever's configuration. More details here:
// https://github.com/golang/go/blob/master/src/net/http/server.go#L2545-L2561
func (params *serverParams) Close(error) error {
	if err := closeServer(params.options.Logger, params.secureServer, params.options.Name, "HTTPS"); err != nil {
		params.err = append(params.err, err)
	}
	if err := closeServer(params.options.Logger, params.insecureServer, params.options.Name, "HTTP"); err != nil {
		params.err = append(params.err, err)
	}
	return utilerrors.NewAggregate(params.err)
}

func closeServer(log logger, svr *http.Server, name, scheme string) error {
	if svr == nil {
		return nil
	}
	if err := svr.Shutdown(context.Background()); err != nil {
		log.Error(fmt.Sprintf("error gracefully closing %s %s server: %v", name, scheme, err))
		return err
	}
	log.Info(fmt.Sprintf("gracefully closed %s %s server", name, scheme))
	return nil
}

func getTLSConfigFromReader(r io.Reader) (*tls.Config, error) {
	if r == nil {
		return nil, errors.New("reader is nil")
	}
	caCert, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	caCertPool := x509.NewCertPool()
	if ok := caCertPool.AppendCertsFromPEM(caCert); !ok {
		return nil, errors.New("could not append cert to pool")
	}
	tlsConfig := &tls.Config{
		ClientCAs:  caCertPool,
		ClientAuth: tls.RequireAndVerifyClientCert,
	}
	tlsConfig.BuildNameToCertificate()
	return tlsConfig, nil
}
