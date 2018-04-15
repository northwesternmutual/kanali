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
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"regexp"
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

func TestLoadPublicKey(t *testing.T) {
	_, err := loadPublicKey(nil)
	assert.Error(t, err)

	_, err = loadPublicKey(bytes.NewReader([]byte("foo")))
	assert.Error(t, err)

	_, err = loadPublicKey(new(errorReader))
	assert.Error(t, err)

	key, _ := rsa.GenerateKey(rand.Reader, 1024)
	data, _ := x509.MarshalPKIXPublicKey(&key.PublicKey)

	result, err := loadPublicKey(bytes.NewReader(pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: data,
	})))
	assert.Nil(t, err)
	assert.Equal(t, &key.PublicKey, result)

	_, err = loadPublicKey(bytes.NewReader(pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(&key.PublicKey),
	})))
	assert.NotNil(t, err)
}

func TestGenerateRandomBytes(t *testing.T) {
	tests := []struct {
		size int
		err  bool
	}{
		{
			size: 1,
		},
		{
			size: 10,
		},
		{
			size: 0,
			err:  true,
		},
	}

	for _, test := range tests {
		data, err := GenerateRandomBytes(test.size)
		if test.err {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
			assert.True(t, regexp.MustCompile(fmt.Sprintf("[0-9a-zA-Z]{%d}", test.size)).Match(data))
		}
	}
}
