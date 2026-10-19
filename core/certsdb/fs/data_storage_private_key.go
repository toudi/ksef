package fs

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"ksef/core/keyring"
	"ksef/internal/utils"
	"os"
)

func (d *FSDataStorage) privateKeyFileName(keyId string) string {
	return d.fileName(keyId, "-pkey.pem")
}

func (d *FSDataStorage) GetPrivateKeyBytes(ctx context.Context, keyId string) ([]byte, error) {
	rawBytes, err := os.ReadFile(d.privateKeyFileName(keyId))
	if err != nil {
		return nil, err
	}

	if bytes.HasPrefix(rawBytes, encryptionHeader) {
		var decryptionKey []byte
		decryptionKey, err = d.loadEncryptionKey(keyId)
		if err != nil {
			return nil, err
		}
		rawBytes, err = utils.GCMAESDecrypt(rawBytes[len(encryptionHeader):], decryptionKey)
	}

	return rawBytes, err
}

// StorePrivateKeyBytes creates-or-overwrites the private key (PEM) under the given
// key ID, encrypting it with a fresh key stored in the keyring. In the
// filesystem implementation the key ID is the certificate UID.
func (d *FSDataStorage) StorePrivateKeyBytes(ctx context.Context, keyId string, pemBytes []byte) error {
	return d.encryptPrivateKeyBytesToFile(keyId, pemBytes, d.privateKeyFileName(keyId))
}

func (d *FSDataStorage) loadEncryptionKey(keyId string) ([]byte, error) {
	var decryptionKey string
	var err error

	if decryptionKey, err = d.kr.Get(
		keyring.AppPrefix,
		keyring.EncryptionKeyUser,
		keyring.PrivateKeyEncryptionKey(keyId),
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
