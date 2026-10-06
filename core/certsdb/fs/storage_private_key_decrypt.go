package fs

import (
	"context"
	"ksef/core/keyring"
	"ksef/internal/logging"
	"log/slog"
	"os"
)

func (f *fsCertsStorage) DecryptPrivateKey(ctx context.Context, certId string) error {
	cert, err := f.getCertByID(certId)
	if err != nil {
		return err
	}

	privateKeyBytes, err := f.GetPrivateKeyBytes(ctx, cert.UID)
	if err != nil {
		return err
	}
	if err := os.WriteFile(f.privateKeyFilename(cert), privateKeyBytes, 0600); err != nil {
		logging.CertsDBLogger.Error("błąd zapisu pliku", slog.Any("err", err))
	}
	return f.kr.Delete(keyring.AppPrefix, keyring.EncryptionKeyUser, keyring.PrivateKeyEncryptionKey(cert.UID))
}
