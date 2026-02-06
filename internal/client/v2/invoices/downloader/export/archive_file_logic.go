package export

import (
	"context"
	"io"
	"ksef/internal/http"
	"os"
	"path/filepath"
)

const archiveFilename = "export.zip"

func (ah *archiveHandler) downloadPart(ctx context.Context, part exportStatusPart, client *http.Client) error {
	destFile, err := os.Create(filepath.Join(ah.workDir, part.PartName))
	if err != nil {
		return err
	}
	defer destFile.Close()
	return client.Download(
		ctx, part.URL, destFile,
	)
}

func (ah *archiveHandler) decryptPart(part exportStatusPart) error {
	encryptedData, err := os.ReadFile(filepath.Join(ah.workDir, part.PartName))
	if err != nil {
		return err
	}
	decryptedBytes, err := ah.cipher.Decrypt(encryptedData, true)
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(ah.workDir, part.decryptedFilename()), decryptedBytes, 0644)
}

func (ah *archiveHandler) concatenateParts(parts []exportStatusPart) error {
	outFile, err := os.Create(filepath.Join(ah.workDir, archiveFilename))
	if err != nil {
		return err
	}
	defer outFile.Close()

	buffer := make([]byte, 10*1024*1024) // 10 mb buffer

	for _, part := range parts {
		// Open each input file
		inFile, err := os.Open(filepath.Join(ah.workDir, part.decryptedFilename()))
		if err != nil {
			return err
		}
		defer inFile.Close()

		// Copy the file in chunks to the output
		_, err = io.CopyBuffer(outFile, inFile, buffer)
		if err != nil {
			return err
		}
	}

	return nil
}
