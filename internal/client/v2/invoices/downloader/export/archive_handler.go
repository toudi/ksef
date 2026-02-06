package export

import (
	"context"
	"ksef/internal/encryption"
	"ksef/internal/http"
	"ksef/internal/logging"
	"os"
)

type archiveHandler struct {
	workDir  string
	cipher   *encryption.Cipher
	contents archiveContents
}

func NewExportArchiveHandler(
	cipher *encryption.Cipher,
) (*archiveHandler, error) {
	tmpDir, err := os.MkdirTemp("", "ksef-export-*")
	if err != nil {
		return nil, err
	}

	return &archiveHandler{
		workDir: tmpDir,
		cipher:  cipher,
	}, nil
}

func (ah *archiveHandler) DownloadExportFile(
	ctx context.Context,
	statusResponse *exportStatusResponse,
) error {
	unauthedClient := http.NewClient("")

	for _, part := range statusResponse.Package.Parts {
		logging.DownloadLogger.Info("Rozpoczynam pobieranie paczki", "part", part.OrdinalNumber)
		if err := ah.downloadPart(ctx, part, unauthedClient); err != nil {
			return err
		}
		logging.DownloadLogger.Info("Rozszyfrowuję paczkę", "part", part.OrdinalNumber)
		if err := ah.decryptPart(part); err != nil {
			return err
		}
	}
	logging.DownloadLogger.Info("Wszystkie paczki pobrane i rozszyfrowane. łączę archiwum")
	if err := ah.concatenateParts(statusResponse.Package.Parts); err != nil {
		return err
	}
	if err := ah.loadMetadata(); err != nil {
		return err
	}
	return nil
}

func (ah *archiveHandler) Close() error {
	return os.RemoveAll(ah.workDir)
}
