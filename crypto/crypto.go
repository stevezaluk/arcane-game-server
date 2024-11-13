package crypto

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"

	arcaneErrors "github.com/stevezaluk/arcane-game-server/errors"
)

type KeyPair struct {
	PrivateKey        *rsa.PrivateKey
	PublicKey         rsa.PublicKey
	PublicKeyChecksum string
	PublicKeyPem      string
}

func New() (KeyPair, error) {
	var ret KeyPair

	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return ret, arcaneErrors.ErrKeyGenerationFailed
	}

	err = privateKey.Validate()
	if err != nil {
		return ret, arcaneErrors.ErrKeysNotValid
	}

	ret.PrivateKey = privateKey
	ret.PublicKey = privateKey.PublicKey

	ret.PublicKeyPem = string(PublicKeyToPEM(ret.PublicKey))

	return ret, nil
}

func PublicKeyToPEM(publicKey rsa.PublicKey) []byte {
	marshaledKey := x509.MarshalPKCS1PublicKey(&publicKey)
	block := &pem.Block{Type: "RSA PUBLIC KEY", Bytes: marshaledKey}

	publicKeyBytes := pem.EncodeToMemory(block)

	return publicKeyBytes
}

func DecryptMessage(message string, privateKey *rsa.PrivateKey) (string, error) {
	cipherText, err := base64.StdEncoding.WithPadding(base64.StdPadding).DecodeString(message)
	if err != nil {
		return "", arcaneErrors.ErrBase64DecodeFailed
	}

	opts := &rsa.OAEPOptions{Hash: crypto.SHA256}

	plainText, err := privateKey.Decrypt(nil, cipherText, opts)
	if err != nil {
		return "", arcaneErrors.ErrDecryptionFailed
	}

	ret := string(plainText)

	return ret, nil
}
