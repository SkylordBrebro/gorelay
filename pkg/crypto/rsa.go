package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
)

// Default RSA public key used by the game server
const defaultPublicKey = `-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDCKFctVrhfF3m2Kes0FBL/JFeO
cmNg9eJz8k/hQy1kadD+XFUpluRqa//Uxp2s9W2qE0EoUCu59ugcf/p7lGuL99Uo
SGmQEynkBvZct+/M40L0E0rZ4BVgzLOJmIbXMp0J4PnPcb6VLZvxazGcmSfjauC7
F3yWYqUbZd/HCBtawwIDAQAB
-----END PUBLIC KEY-----`

var (
	publicKey *rsa.PublicKey
)

func init() {
	var err error
	publicKey, err = parsePublicKey(defaultPublicKey)
	if err != nil {
		panic(fmt.Sprintf("failed to parse public key: %v", err))
	}
}

// parsePublicKey parses a PEM encoded RSA public key
func parsePublicKey(pemStr string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pemStr))
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %v", err)
	}

	switch pub := pub.(type) {
	case *rsa.PublicKey:
		return pub, nil
	default:
		return nil, fmt.Errorf("unsupported public key type")
	}
}

// Encrypt encrypts data using RSA PKCS1v15
func Encrypt(data string) string {
	encrypted, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, []byte(data))
	if err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(encrypted)
}

// SetPublicKey sets a custom RSA public key
func SetPublicKey(pemStr string) error {
	key, err := parsePublicKey(pemStr)
	if err != nil {
		return err
	}
	publicKey = key
	return nil
}
