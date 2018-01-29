package server

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"

	"golang.org/x/sync/errgroup"
	"k8s.io/apimachinery/pkg/util/errors"
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

// PrepareServer will construct net.Listerer instantiations for the requested
// server(s). If an error is encounted, it will be hidden within the returned
// type for evaluation by that type's methods.
func PrepareServer(opts *Options) *serverParams {
	params := &serverParams{options: opts}
	insecureAddr, secureAddr := fmt.Sprintf("%s:%d", opts.InsecureAddr, opts.InsecurePort), fmt.Sprintf("%s:%d", opts.SecureAddr, opts.SecurePort)

	if opts.SecurePort > 0 && len(opts.TLSCa) > 0 { // client has requested an HTTPS server with mutual TLS
		tlsConfig, err := opts.getTLSConfig()
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
func (params *serverParams) Run() error {
	if params.err != nil && len(params.err) > 0 {
		return errors.NewAggregate(params.err)
	}

	var g errgroup.Group
	if params.secureServer != nil {
		params.options.Logger.Info(fmt.Sprintf("starting %s %s server on %s:%d", params.options.Name, "HTTP", params.options.SecureAddr, params.options.SecurePort))
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
func (params *serverParams) Close() error {
	if err := closeServer(params.options.Logger, params.secureServer, params.options.Name, "HTTPS"); err != nil {
		params.err = append(params.err, err)
	}
	if err := closeServer(params.options.Logger, params.insecureServer, params.options.Name, "HTTP"); err != nil {
		params.err = append(params.err, err)
	}
	return errors.NewAggregate(params.err)
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

func (opts *Options) getTLSConfig() (*tls.Config, error) {
	caCert, err := ioutil.ReadFile(opts.TLSCa)
	if err != nil {
		return nil, err
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	tlsConfig := &tls.Config{
		ClientCAs:  caCertPool,
		ClientAuth: tls.RequireAndVerifyClientCert,
	}
	tlsConfig.BuildNameToCertificate()
	return tlsConfig, nil
}
