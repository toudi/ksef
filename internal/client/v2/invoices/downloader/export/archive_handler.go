package export

import (
	"context"
	"ksef/internal/encryption"
	"ksef/internal/http"
	"log/slog"
	"os"
)

type archiveHandler struct {
	workDir  string
	cipher   *encryption.Cipher
	contents *archiveContents
	logger   *slog.Logger
}

func NewExportArchiveHandler(
	cipher *encryption.Cipher,
	logger *slog.Logger,
) (*archiveHandler, error) {
	tmpDir, err := os.MkdirTemp("", "ksef-export-*")
	if err != nil {
		return nil, err
	}

	logger.Debug("katalog roboczy dla pobranego eksportu", "dir", tmpDir)

	return &archiveHandler{
		workDir: tmpDir,
		cipher:  cipher,
		logger:  logger,
	}, nil
}

func (ah *archiveHandler) DownloadExportFile(
	ctx context.Context,
	statusResponse *exportStatusResponse,
) error {
	unauthedClient := http.NewClient("")

	for _, part := range statusResponse.Package.Parts {
		ah.logger.Info("Rozpoczynam pobieranie paczki", "part", part.OrdinalNumber)
		if err := ah.downloadPart(ctx, part, unauthedClient); err != nil {
			return err
		}
		ah.logger.Info("Rozszyfrowuję paczkę", "part", part.OrdinalNumber)
		if err := ah.decryptPart(part); err != nil {
			return err
		}
	}
	ah.logger.Info("Wszystkie paczki pobrane i rozszyfrowane. łączę archiwum")
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
	ah.logger.Debug("Zamykam archiwum ZIP")
	if err := ah.contents.zipReader.Close(); err != nil {
		return err
	}
	ah.logger.Debug("Usuwam katalog roboczy", "dir", ah.workDir)
	return os.RemoveAll(ah.workDir)
}
