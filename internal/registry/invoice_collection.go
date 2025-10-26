package registry

import (
	"errors"
	"fmt"
	"ksef/internal/logging"
	"os"
	"path/filepath"
	"strings"
)

var (
	ErrAlreadySynced        = errors.New("all invoices already sent")
	ErrUnableToDetectIssuer = errors.New("unable to detect issuer")
)

// v2 api specs wants us to group invoices by their form code so we have to keep them in hash
type InvoiceFormCode struct {
	SystemCode    string `xml:"kodSystemowy,attr" json:"systemCode"`
	SchemaVersion string `xml:"wersjaSchemy,attr" json:"schemaVersion"`
	Value         string `xml:",chardata" json:"value"`
}

type CollectionFile struct {
	Filename string
	Checksum string
}

type InvoiceCollection struct {
	Issuer string
	Files  map[InvoiceFormCode][]CollectionFile
}

func (r *InvoiceRegistry) InvoiceCollection() (*InvoiceCollection, error) {
	if r.collection != nil {
		return r.collection, nil
	}

	files, err := os.ReadDir(r.Dir)

	if err != nil {
		return nil, fmt.Errorf("cannot read list of files from %s: %v", r.sourcePath, err)
	}

	collection := &InvoiceCollection{
		Files: make(map[InvoiceFormCode][]CollectionFile),
	}

	var fileName string
	var fullFileName string
	var parsedInvoice *XMLInvoice
	var checksum string
	var allFiles int
	var syncedFiles int

	for _, file := range files {
		fileName = file.Name()
		fullFileName = filepath.Join(r.Dir, fileName)

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
			invoice, _ := r.GetInvoiceByChecksum(checksum)
			if invoice.KSeFReferenceNumber != "" {
				logging.UploadLogger.Info(
					"invoice was already uploaded - skipping",
					"invoice",
					fullFileName,
				)
				syncedFiles += 1
				continue
			}
			collection.Files[parsedInvoice.HeaderFormCode] = append(collection.Files[parsedInvoice.HeaderFormCode], CollectionFile{
				Filename: fullFileName, Checksum: checksum,
			})

			_, err = r.Upsert(Invoice{
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
		return nil, ErrUnableToDetectIssuer
	}

	r.collection = collection

	return collection, nil
}
