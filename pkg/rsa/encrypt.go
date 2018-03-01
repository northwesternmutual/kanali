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
)

type option func(*config)

type config struct {
	base64Encode    bool
	encryptionLabel []byte
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

	if !cfg.base64Encode {
		return data, nil
	}

	var b64 []byte
	base64.StdEncoding.Encode(b64, data)
	return b64, nil
}

func Base64Encode() option {
	return option(func(cfg *config) {
		cfg.base64Encode = true
	})
}

func WithEncryptionLabel(label []byte) option {
	return option(func(cfg *config) {
		cfg.encryptionLabel = label
	})
}
