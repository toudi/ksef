package export

import (
	"bytes"
	"encoding/json"
	"io"
	"ksef/internal/client/v2/types/invoices"
	"os"
	"path/filepath"
)

const metadataFilename = "_metadata.json"

type archiveContents struct {
	Invoices []invoices.InvoiceMetadata `json:"invoices"`
}

func (ah *archiveHandler) ReadInvoice(ksefRefNo string, dest *bytes.Buffer) error {
	dest.Reset()

	invoiceFile, err := os.Open(filepath.Join(ah.workDir, ksefRefNo+".xml"))
	if err != nil {
		return err
	}
	defer invoiceFile.Close()

	if _, err = io.Copy(dest, invoiceFile); err != nil {
		return err
	}

	return nil
}

func (ah *archiveHandler) loadMetadata() error {
	metadataFile, err := os.Open(filepath.Join(ah.workDir, metadataFilename))
	if err != nil {
		return err
	}
	defer metadataFile.Close()

	return json.NewDecoder(metadataFile).Decode(&ah.contents)
}
