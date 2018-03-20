package rsa

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type errorReader struct{}

func (r *errorReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("foo bar car")
}

func TestLoadDecryptionKey(t *testing.T) {
	_, err := loadDecryptionKey(nil)
	assert.Error(t, err)

	_, err = loadDecryptionKey(bytes.NewReader([]byte("foo")))
	assert.Error(t, err)

	_, err = loadDecryptionKey(new(errorReader))
	assert.Error(t, err)

	key, _ := rsa.GenerateKey(rand.Reader, 1024)
	keyData := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	})

	result, err := loadDecryptionKey(bytes.NewReader(keyData))
	assert.Nil(t, err)
	assert.Equal(t, key, result)
}
