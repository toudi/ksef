package encryption

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

func EncryptMessageWithCertificate(certificateFileBytes []byte, message []byte) ([]byte, error) {
	block, _ := pem.Decode(certificateFileBytes)
	parsedCert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("could not parse certificate: %w", err)
	}
	var publicKey *rsa.PublicKey
	var ok bool
	if publicKey, ok = parsedCert.PublicKey.(*rsa.PublicKey); !ok {
		return nil, fmt.Errorf("cannot parse public key: %v", err)
	}
	encryptedBytes, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, message, nil)
	if err != nil {
		return nil, fmt.Errorf("could not encrypt cipher key with finance ministry's public key: %v", err)
	}

	return encryptedBytes, nil
}
