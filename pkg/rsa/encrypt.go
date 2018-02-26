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

func Encrypt(data string, key *rsa.PublicKey, options ...option) ([]byte, error) {
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
