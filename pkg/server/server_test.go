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
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/northwesternmutual/kanali/pkg/log"
	"github.com/northwesternmutual/kanali/test/builder"
)

type errorReader struct{}

func (r *errorReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("foo bar car")
}

func TestGetTLSConfigFromReader(t *testing.T) {
	// mimic a cert template
	tmpl := builder.CreateGenericCertificateTemplate()
	tmpl.IsCA = true
	tmpl.KeyUsage = x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign
	// create a private key
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.Nil(t, err)
	// create a certificate in DER format
	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &key.PublicKey, key)
	assert.Nil(t, err)
	// convert certificate to PEM format
	pem := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})

	cfg, err := getTLSConfigFromReader(bytes.NewReader(pem))
	assert.Nil(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, 1, len(cfg.ClientCAs.Subjects()))

	_, err = getTLSConfigFromReader(bytes.NewReader([]byte("foo")))
	assert.NotNil(t, err)

	_, err = getTLSConfigFromReader(nil)
	assert.NotNil(t, err)

	_, err = getTLSConfigFromReader(new(errorReader))
	assert.NotNil(t, err)
}

func TestPrepare(t *testing.T) {
	opts := &Options{
		Name:         "test",
		InsecureAddr: "1.2.3.4",
		SecureAddr:   "4.3.2.1",
		InsecurePort: 1234,
		SecurePort:   4321,
		TLSKey:       "/foo/tls.key",
		TLSCert:      "/foo/tls.crt",
		TLSCa:        "/foo/tls.ca",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}),
		Logger: log.WithContext(nil).Sugar(),
	}

	// mimic a cert template
	tmpl := builder.CreateGenericCertificateTemplate()
	tmpl.IsCA = true
	tmpl.BasicConstraintsValid = true
	tmpl.KeyUsage = x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign
	// create a private key
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.Nil(t, err)
	// create a certificate in DER format
	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &key.PublicKey, key)
	assert.Nil(t, err)
	// convert certificate to PEM format
	pem := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})

	params := prepare(opts, bytes.NewReader(pem))
	assert.NotNil(t, params)
	assert.Nil(t, params.err)
	assert.NotNil(t, params.secureServer)
	assert.NotNil(t, params.insecureServer)
	assert.Equal(t, params.options, opts)
	assert.Equal(t, params.secureServer.Addr, "4.3.2.1:4321")
	assert.Equal(t, params.insecureServer.Addr, "1.2.3.4:1234")

	params = prepare(opts, nil)
	assert.NotNil(t, params)
	assert.NotNil(t, params.err)
	assert.Equal(t, params.err[0].Error(), "reader is nil")
	assert.NotNil(t, params.insecureServer)

	opts.TLSCa = ""
	params = prepare(opts, bytes.NewReader(pem))
	assert.NotNil(t, params)
	assert.Nil(t, params.err)
	assert.NotNil(t, params.secureServer)
	assert.NotNil(t, params.insecureServer)
}

func TestIsEnabled(t *testing.T) {
	assert.False(t, (&serverParams{
		err: []error{errors.New("foo")},
	}).IsEnabled())
	assert.True(t, (&serverParams{
		options: &Options{
			InsecurePort: 1,
		},
	}).IsEnabled())
}

func TestRun(t *testing.T) {
	params := serverParams{
		err: []error{errors.New("foo")},
	}
	err := params.Run(context.Background())
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "foo")
}
