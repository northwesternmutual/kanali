package tls

import "crypto/x509"

// VerifyPeerCertificate will validate that the server cert is signed from a
// certificate in a given pool. Common name and SANS validation is not performed.
func VerifyPeerCertificate(
	pool *x509.CertPool,
) func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
	return func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
		opts := x509.VerifyOptions{
			Roots: pool,
		}
		cert, err := x509.ParseCertificate(rawCerts[0])
		if err != nil {
			return err
		}
		_, err = cert.Verify(opts)
		return err
	}
}
