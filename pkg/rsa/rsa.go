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
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
)

type option func(*config)

type config struct {
	base64Encode, base64Decode bool
	encryptionLabel            []byte
}

var EncryptionLabel = []byte("kanali")

func Encrypt(data []byte, key *rsa.PublicKey, options ...option) ([]byte, error) {
	cfg := new(config)
	for _, op := range options {
		op(cfg)
	}

	data, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, key, data, cfg.encryptionLabel)
	if err != nil {
		return nil, err
	}

	h := make([]byte, hex.EncodedLen(len(data)))
	data = h[:hex.Encode(h, data)]

	if !cfg.base64Encode {
		return data, nil
	}

	b64 := make([]byte, base64.StdEncoding.EncodedLen(len(data)))
	base64.StdEncoding.Encode(b64, data)
	return b64, nil
}

func Decrypt(cipherText []byte, key *rsa.PrivateKey, options ...option) ([]byte, error) {
	cfg := new(config)
	for _, op := range options {
		op(cfg)
	}

	if cfg.base64Decode {
		b64 := make([]byte, base64.StdEncoding.DecodedLen(len(cipherText)))
		n, err := base64.StdEncoding.Decode(b64, cipherText)
		if err != nil {
			return nil, err
		}
		cipherText = b64[:n]
	}

	hexDecodedCipherText := make([]byte, hex.DecodedLen(len(cipherText)))
	n, err := hex.Decode(hexDecodedCipherText, cipherText)
	if err != nil {
		return nil, err
	}

	return rsa.DecryptOAEP(sha256.New(), rand.Reader, key, hexDecodedCipherText[:n], cfg.encryptionLabel)
}

func Base64Encode() option {
	return option(func(cfg *config) {
		cfg.base64Encode = true
	})
}

func Base64Decode() option {
	return option(func(cfg *config) {
		cfg.base64Decode = true
	})
}

func WithEncryptionLabel(label []byte) option {
	return option(func(cfg *config) {
		cfg.encryptionLabel = label
	})
}
