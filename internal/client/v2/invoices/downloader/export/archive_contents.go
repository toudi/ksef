package export

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"io"
	"ksef/internal/client/v2/types/invoices"
	"ksef/internal/logging"
	"path/filepath"
)

const metadataFilename = "_metadata.json"

type archiveContents struct {
	Invoices  []invoices.InvoiceMetadata `json:"invoices"`
	zipReader *zip.ReadCloser
}

func NewArchiveContents(dirName string) (*archiveContents, error) {
	logging.DownloadLogger.Debug("Otwieram archiwum ZIP")
	zipReader, err := zip.OpenReader(filepath.Join(dirName, archiveFilename))
	if err != nil {
		return nil, err
	}

	logging.DownloadLogger.Debug("Odczytuję plik metadanych", "path", filepath.Join(dirName, archiveFilename)+" :: "+metadataFilename)
	metadataFile, err := zipReader.Open(metadataFilename)
	if err != nil {
		return nil, err
	}
	defer metadataFile.Close()

	contents := &archiveContents{
		zipReader: zipReader,
	}

	logging.DownloadLogger.Debug("Dekoduję tablicę faktur do pamięci")
	if err = json.NewDecoder(metadataFile).Decode(contents); err != nil {
		return nil, err
	}

	logging.DownloadLogger.Debug("Wczytywanie metadanych zakończone")
	return contents, nil
}

func (ah *archiveHandler) ReadInvoice(ksefRefNo string, dest *bytes.Buffer) error {
	dest.Reset()

	invoiceFile, err := ah.contents.zipReader.Open(ksefRefNo + ".xml")
	if err != nil {
		return err
	}
	defer invoiceFile.Close()

	if _, err = io.Copy(dest, invoiceFile); err != nil {
		return err
	}

	return nil
}

func (ac *archiveContents) Close() error {
	return ac.zipReader.Close()
}
