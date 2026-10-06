package fs

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"ksef/core/keyring"
	"ksef/internal/utils"
	"os"
	"path/filepath"
)

func (f *fsCertsStorage) privateKeyFilename(cert *fsCertificate) string {
	return filepath.Join(filepath.Dir(certificatesDBFile), string(cert.EnvironmentId)+"-"+cert.UID+"-pkey.pem")
}

func (f *fsCertsStorage) GetPrivateKeyBytes(ctx context.Context, certId string) ([]byte, error) {
	cert, err := f.getCertByID(certId)
	if err != nil {
		return nil, err
	}

	rawBytes, err := os.ReadFile(f.privateKeyFilename(cert))
	if err != nil {
		return nil, err
	}

	if bytes.HasPrefix(rawBytes, encryptionHeader) {
		var decryptionKey []byte
		decryptionKey, err = f.loadEncryptionKey(cert.UID)
		if err != nil {
			return nil, err
		}
		rawBytes, err = utils.GCMAESDecrypt(rawBytes[len(encryptionHeader):], decryptionKey)
	}

	return rawBytes, err
}

func (f *fsCertsStorage) loadEncryptionKey(certId string) ([]byte, error) {
	var decryptionKey string
	var err error

	if decryptionKey, err = f.kr.Get(
		keyring.AppPrefix,
		keyring.EncryptionKeyUser,
		keyring.PrivateKeyEncryptionKey(certId),
	); err != nil {
		return nil, err
	}

	if decryptionKey == "" {
		return nil, errors.New("nie znaleziono klucza szyfrującego")
	}
	// decryption key is base64 encoded since we're dealing with bytes but zalando lib stores strings so let's just be
	// safe than sorry
	return base64.StdEncoding.DecodeString(decryptionKey)
}
