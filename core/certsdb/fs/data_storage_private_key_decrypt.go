package fs

import (
	"context"
	"ksef/core/keyring"
	"ksef/internal/logging"
	"log/slog"
	"os"
)

func (d *FSDataStorage) DecryptPrivateKey(ctx context.Context, keyId string) error {
	privateKeyBytes, err := d.GetPrivateKeyBytes(ctx, keyId)
	if err != nil {
		return err
	}
	if err := os.WriteFile(d.privateKeyFileName(keyId), privateKeyBytes, 0600); err != nil {
		logging.CertsDBLogger.Error("błąd zapisu pliku", slog.Any("err", err))
	}
	return d.kr.Delete(keyring.AppPrefix, keyring.EncryptionKeyUser, keyring.PrivateKeyEncryptionKey(keyId))
}
