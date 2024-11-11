package crypto

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"log/slog"

	"github.com/stevezaluk/arcane-game-server/errors"
)

func GenerateKeyPair() (rsa.PrivateKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return *privateKey, errors.ErrKeyGenerationFailed
	}

	err = privateKey.Validate()
	if err != nil {
		return *privateKey, errors.ErrKeysNotValid
	}

	return *privateKey, nil
}

func PublicKeyToPEM(publicKey rsa.PublicKey) []byte {
	marshaledKey := x509.MarshalPKCS1PublicKey(&publicKey)
	block := &pem.Block{Type: "RSA PUBLIC KEY", Bytes: marshaledKey}

	publicKeyBytes := pem.EncodeToMemory(block)

	return publicKeyBytes
}

func DecryptMessage(message string, privateKey *rsa.PrivateKey) string {
	cipherText, err := base64.StdEncoding.WithPadding(base64.StdPadding).DecodeString(message)
	if err != nil {
		slog.Error("Failed to decrypt base64 encoded cipher text")
		return ""
	}

	opts := &rsa.OAEPOptions{Hash: crypto.SHA256}

	plainText, err := privateKey.Decrypt(nil, cipherText, opts)
	if err != nil {
		slog.Error("Failed to decrypt cipher text provided by the client")
		return ""
	}

	return string(plainText)
}
