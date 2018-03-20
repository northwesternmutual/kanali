package rsa

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io"
	"io/ioutil"
	"os"
)

func LoadDecryptionKey(location string) (*rsa.PrivateKey, error) {
	file, err := os.Open(location)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return loadDecryptionKey(file)
}

func loadDecryptionKey(reader io.Reader) (*rsa.PrivateKey, error) {
	if reader == nil {
		return nil, errors.New("reader is nil")
	}
	// read in private key
	keyBytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	// create a pem block from the private key provided
	block, _ := pem.Decode(keyBytes)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, errors.New("failed to decode PEM block containing rsa private key")
	}
	// parse the pem block into a private key
	return x509.ParsePKCS1PrivateKey(block.Bytes)
}
