package tls

import (
	"crypto/x509"
	"sync"
)

var (
	once sync.Once
	pool *x509.CertPool
	err  error
)

func GetSystemCertPool() (*x509.CertPool, error) {
	once.Do(initSystemCertPool)
	if err != nil {
		return nil, err
	}
	return pool, nil
}

func initSystemCertPool() {
	pool, err = x509.SystemCertPool()
}
