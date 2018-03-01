package kanali

import (
	"crypto/rand"
	"crypto/rsa"
)

var (
	TestApiKeyDecryptionPrivateKey *rsa.PrivateKey
	TestApiKeyDecryptionPublicKey  *rsa.PublicKey
)

func initApiKeyDecryptionKeys(bits int) error {
	priv, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return err
	}
	TestApiKeyDecryptionPrivateKey = priv
	TestApiKeyDecryptionPublicKey = &priv.PublicKey
	return nil
}
