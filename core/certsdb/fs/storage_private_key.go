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
	return f.privateKeyFilenameInternal(cert.EnvironmentId, cert.UID)
}

func (f *fsCertsStorage) privateKeyFilenameInternal(gatewayId string, uid string) string {
	return filepath.Join(filepath.Dir(certificatesDBFile), gatewayId+"-"+uid+"-pkey.pem")
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

func (f *fsCertsStorage) GetPrivateKeyId(ctx context.Context, certId string) (string, error) {
	cert, err := f.getCertByID(certId)
	if err != nil {
		return "", err
	}

	return cert.UID, nil
}

func (f *fsCertsStorage) StorePrivateKey(ctx context.Context, pemBytes []byte, gatewayId string) (string, error) {
	uid, err := f.generateUID()
	if err != nil {
		return "", err
	}
	privateKeyFilename := f.privateKeyFilenameInternal(gatewayId, uid)
	privateKeyFile, err := os.Create(privateKeyFilename)
	if err != nil {
		return "", err
	}
	defer privateKeyFile.Close()
	err = f.encryptPrivateKeyBytes(uid, pemBytes, privateKeyFile)
	return uid, err
}
