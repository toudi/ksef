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

func (d *FSDataStorage) IsPrivateKeyEncrypted(ctx context.Context, keyId string) (bool, error) {
	privateKeyBytes, err := os.ReadFile(d.privateKeyFileName(keyId))
	if err != nil {
		return false, err
	}

	return bytes.HasPrefix(privateKeyBytes, encryptionHeader), nil
}

func (d *FSDataStorage) EncryptPrivateKey(ctx context.Context, keyId string) error {
	privateKeyFilename := d.privateKeyFileName(keyId)
	privateKeyBytes, err := os.ReadFile(privateKeyFilename)
	if err != nil {
		return err
	}
	privateKeyFile, err := os.OpenFile(privateKeyFilename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer privateKeyFile.Close()
	return d.encryptPrivateKeyBytes(keyId, privateKeyBytes, privateKeyFile)
}

func (d *FSDataStorage) encryptPrivateKeyBytes(keyId string, pemBytes []byte, dest io.Writer) error {
	encryptionKey := make([]byte, 16)
	_, err := rand.Read(encryptionKey)
	if err != nil {
		return err
	}
	// let's save the key to keyring first before proceeding further
	if err = d.kr.Set(
		keyring.AppPrefix,
		keyring.EncryptionKeyUser,
		keyring.PrivateKeyEncryptionKey(keyId),
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

func (d *FSDataStorage) encryptPrivateKeyBytesToFile(keyId string, pemBytes []byte, privateKeyFilename string) error {
	privateKeyFile, err := os.OpenFile(privateKeyFilename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer privateKeyFile.Close()
	return d.encryptPrivateKeyBytes(keyId, pemBytes, privateKeyFile)
}
