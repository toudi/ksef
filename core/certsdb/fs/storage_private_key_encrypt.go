package fs

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"io"
	"ksef/core/keyring"
	"ksef/internal/utils"
	"os"
)

var encryptionHeader = []byte{
	'E', 'P', 'K', 0x00, 0x99,
}

func (f *fsCertsStorage) IsPrivateKeyEncrypted(ctx context.Context, certId string) (bool, error) {
	cert, err := f.getCertByID(certId)
	if err != nil {
		return false, err
	}
	privateKeyBytes, err := os.ReadFile(f.privateKeyFilename(cert))
	if err != nil {
		return false, err
	}

	return bytes.HasPrefix(privateKeyBytes, encryptionHeader), nil
}

func (f *fsCertsStorage) EncryptPrivateKey(ctx context.Context, certId string) error {
	cert, err := f.getCertByID(certId)
	if err != nil {
		return err
	}

	privateKeyFilename := f.privateKeyFilename(cert)
	privateKeyBytes, err := os.ReadFile(privateKeyFilename)
	if err != nil {
		return err
	}
	privateKeyFile, err := os.OpenFile(privateKeyFilename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer privateKeyFile.Close()
	return f.encryptPrivateKeyBytes(cert.UID, privateKeyBytes, privateKeyFile)
}

func (f *fsCertsStorage) encryptPrivateKeyBytes(certId string, pemBytes []byte, dest io.Writer) error {
	encryptionKey := make([]byte, 16)
	_, err := rand.Read(encryptionKey)
	if err != nil {
		return err
	}
	// let's save the key to keyring first before proceeding further
	if err = f.kr.Set(
		keyring.AppPrefix,
		keyring.EncryptionKeyUser,
		keyring.PrivateKeyEncryptionKey(certId),
		base64.StdEncoding.EncodeToString(encryptionKey),
	); err != nil {
		return err
	}
	encryptedKeyBytes, err := utils.GCMAESEncrypt(pemBytes, encryptionKey)
	if err != nil {
		return err
	}
	encryptedKeyBytes = bytes.Join([][]byte{encryptionHeader, encryptedKeyBytes}, nil)
	_, err = io.Copy(dest, bytes.NewReader(encryptedKeyBytes))
	return err
}

func (f *fsCertsStorage) encryptPrivateKeyBytesToFile(certId string, pemBytes []byte, privateKeyFilename string) error {
	privateKeyFile, err := os.OpenFile(privateKeyFilename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer privateKeyFile.Close()
	return f.encryptPrivateKeyBytes(certId, pemBytes, privateKeyFile)
}
