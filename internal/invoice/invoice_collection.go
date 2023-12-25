package invoice

import (
	"errors"
	"fmt"

	"ksef/internal/logging"
	registryPkg "ksef/internal/registry"
	"os"
	"path/filepath"
	"strings"
)

type collectionFile struct {
	Filename string
	Checksum string
}

type invoiceCollection struct {
	Issuer string
	Files  []collectionFile
}

var ErrAlreadySynced = errors.New("all invoices already sent")

func InvoiceCollection(
	sourcePath string,
	registry *registryPkg.InvoiceRegistry,
) (*invoiceCollection, error) {
	files, err := os.ReadDir(sourcePath)
	if err != nil {
		return nil, fmt.Errorf("cannot read list of files from %s: %v", sourcePath, err)
	}

	collection := &invoiceCollection{
		Files: make([]collectionFile, 0),
	}

	var fileName string
	var fullFileName string
	var parsedInvoice *XMLInvoice
	var checksum string
	var allFiles int
	var syncedFiles int

	for _, file := range files {
		fileName = file.Name()
		fullFileName = filepath.Join(sourcePath, fileName)

		// do not bother to check if the file is not an XML
		if strings.ToLower(filepath.Ext(fileName)) != ".xml" {
			continue
		}

		// let's use the fact that we know how to parse the issuer to detect
		// whether this is actually a SEI / KSeF invoice file or it happens to
		// be just some random XML file
		if parsedInvoice, _ = ParseInvoice(fullFileName); parsedInvoice != nil {
			allFiles += 1
			// calculate checksum for this invoice.
			checksum, err = hashSourceFile(fullFileName)
			if err != nil {
				return nil, fmt.Errorf("unable to hash source file: %v", err)
			}
			invoice := registry.GetInvoiceByChecksum(checksum)
			if invoice.SEIReferenceNumber != "" {
				logging.UploadLogger.Info(
					"invoice was already uploaded - skipping",
					"invoice",
					fullFileName,
				)
				syncedFiles += 1
				continue
			}
			collection.Files = append(collection.Files, collectionFile{
				Filename: fullFileName, Checksum: checksum,
			})

			_, err = registry.Upsert(registryPkg.Invoice{
				Checksum:        checksum,
				ReferenceNumber: parsedInvoice.InvoiceNumber,
			})

			if err != nil {
				return nil, fmt.Errorf("cannot upsert invoice in registry: %v", err)
			}

			if collection.Issuer == "" {
				collection.Issuer = parsedInvoice.Issuer
			}
		}
	}

	if allFiles == syncedFiles {
		return nil, ErrAlreadySynced
	}

	if collection.Issuer == "" {
		return nil, fmt.Errorf("no issuer found - this should be impossible")
	}

	return collection, nil

}
