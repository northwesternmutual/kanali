package rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRSA(t *testing.T) {
	priv, err := rsa.GenerateKey(rand.Reader, 4096)
	assert.Nil(t, err)

	originalValue := []byte("abc123")
	cipherText, err := Encrypt(originalValue, &priv.PublicKey)
	assert.Nil(t, err)

	plainText, err := Decrypt(cipherText, priv)
	assert.Nil(t, err)
	assert.Equal(t, plainText, originalValue)

	cipherText, err = Encrypt(originalValue, &priv.PublicKey, Base64Encode())
	assert.Nil(t, err)

	plainText, err = Decrypt(cipherText, priv, Base64Decode())
	assert.Nil(t, err)
	assert.Equal(t, plainText, originalValue)
}
