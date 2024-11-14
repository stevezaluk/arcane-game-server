package crypto

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
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
	ret.PublicKeyChecksum = PublicKeyToChecksum(ret.PublicKeyPem)

	return ret, nil
}

func PublicKeyToPEM(publicKey rsa.PublicKey) []byte {
	marshaledKey := x509.MarshalPKCS1PublicKey(&publicKey)
	block := &pem.Block{Type: "RSA PUBLIC KEY", Bytes: marshaledKey}

	publicKeyBytes := pem.EncodeToMemory(block)

	return publicKeyBytes
}

func PEMToPublicKey(pemKey string) (rsa.PublicKey, error) {
	publicKeyBlock, _ := pem.Decode([]byte(pemKey))
	pubKey, err := x509.ParsePKCS1PublicKey(publicKeyBlock.Bytes)
	if err != nil {
		return *pubKey, arcaneErrors.ErrParsePubKeyFailed
	}

	return *pubKey, nil
}

func PublicKeyToChecksum(pemKey string) string {
	hash := sha256.Sum256([]byte(pemKey))
	hashStr := hex.EncodeToString(hash[:])

	return hashStr
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
