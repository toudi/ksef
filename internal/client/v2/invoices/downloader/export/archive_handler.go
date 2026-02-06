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
	contents *archiveContents
}

func NewExportArchiveHandler(
	cipher *encryption.Cipher,
) (*archiveHandler, error) {
	tmpDir, err := os.MkdirTemp("", "ksef-export-*")
	if err != nil {
		return nil, err
	}

	logging.DownloadLogger.Debug("katalog roboczy dla pobranego eksportu", "dir", tmpDir)

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
	var err error

	if ah.contents, err = NewArchiveContents(ah.workDir); err != nil {
		return err
	}
	return nil
}

func (ah *archiveHandler) Close() error {
	logging.DownloadLogger.Debug("Zamykam archiwum ZIP")
	if err := ah.contents.zipReader.Close(); err != nil {
		return err
	}
	logging.DownloadLogger.Debug("Usuwam katalog roboczy", "dir", ah.workDir)
	return os.RemoveAll(ah.workDir)
}
