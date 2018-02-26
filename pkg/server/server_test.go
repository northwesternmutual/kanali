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
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/northwesternmutual/kanali/pkg/log"
	"github.com/northwesternmutual/kanali/test/builder"
)

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
}
