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
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}
}
