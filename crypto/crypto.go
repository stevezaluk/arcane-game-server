package crypto

import (
	"crypto/rand"
	"crypto/rsa"
)

func GenerateKeyPair() (rsa.PrivateKey, rsa.PublicKey) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		panic(err) // panicing here as this is a fatal error
	}

	publicKey := privateKey.PublicKey

	return *privateKey, publicKey
}
