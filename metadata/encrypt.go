package metadata

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"ksef/common/aes"
)

type EncryptedArchive struct {
	cipher *aes.Cipher
	size   int
	hash   []byte
}

func (e *EncryptedArchive) encryptKeyWithCertificate(certificateFile string) ([]byte, error) {
	certFileBytes, err := ioutil.ReadFile(certificateFile)
	if err != nil {
		return nil, fmt.Errorf("cannot read certificate file: %v", err)
	}

	block, _ := pem.Decode(certFileBytes)
	parsedKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("could not parse public key from %s: %v", certificateFile, err)
	}
	var publicKey *rsa.PublicKey
	var ok bool
	if publicKey, ok = parsedKey.(*rsa.PublicKey); !ok {
		return nil, fmt.Errorf("cannot parse public key: %v", err)
	}
	encryptedKeyBytes, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, e.cipher.Key)
	if err != nil {
		return nil, fmt.Errorf("could not encrypt cipher key with finance ministry's public key: %v", err)
	}

	return encryptedKeyBytes, nil
}
