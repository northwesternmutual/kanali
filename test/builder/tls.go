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

package builder

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"math/big"
	"net"
	"time"
)

type tlsBuilder struct {
	curr tlsAssets
}

type tlsAssets struct {
	CACert, CAKey, ServerCert, ServerKey []byte
}

func NewTLSBuilder(dnsNames []string, ips []net.IP) *tlsBuilder {
	b := new(tlsBuilder)

	caCertTmpl := CreateGenericCertificateTemplate()
	caCertTmpl.IsCA = true
	caCertTmpl.KeyUsage = x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign

	caKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	caDER, _ := x509.CreateCertificate(rand.Reader, caCertTmpl, caCertTmpl, &caKey.PublicKey, caKey)
	b.curr.CACert = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: caDER})
	b.curr.CAKey = pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(caKey)})

	caCertTmpl, _ = x509.ParseCertificate(caDER)

	serverCertTmpl := CreateGenericCertificateTemplate()
	serverCertTmpl.KeyUsage = x509.KeyUsageDigitalSignature
	serverCertTmpl.DNSNames = dnsNames
	serverCertTmpl.IPAddresses = ips

	serverKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	serverDER, _ := x509.CreateCertificate(rand.Reader, serverCertTmpl, caCertTmpl, &serverKey.PublicKey, caKey)
	b.curr.ServerCert = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: serverDER})
	b.curr.ServerKey = pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(serverKey)})

	return b
}

func (b *tlsBuilder) NewOrDie() *tlsAssets {
	return &b.curr
}

func CreateGenericCertificateTemplate() *x509.Certificate {
	return &x509.Certificate{
		SerialNumber:          big.NewInt(1653),
		NotBefore:             time.Now().AddDate(0, 0, -1),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}
}
