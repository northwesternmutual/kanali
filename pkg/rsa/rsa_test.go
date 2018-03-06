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

	plainText, err = Decrypt(cipherText, priv, Base64Decode())
	assert.NotNil(t, err)

	cipherText, err = Encrypt(originalValue, &priv.PublicKey, Base64Encode(), WithEncryptionLabel(EncryptionLabel))
	assert.Nil(t, err)

	plainText, err = Decrypt(cipherText, priv, Base64Decode(), WithEncryptionLabel(EncryptionLabel))
	assert.Nil(t, err)
	assert.Equal(t, plainText, originalValue)

	plainText, err = Decrypt(cipherText, priv, Base64Decode())
	assert.NotNil(t, err)

	plainText, err = Decrypt(cipherText, priv, WithEncryptionLabel(EncryptionLabel))
	assert.NotNil(t, err)
}
