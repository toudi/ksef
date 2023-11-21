package common

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type invoiceCollection struct {
	Issuer string
	Files  []string
}

func InvoiceCollection(sourcePath string) (*invoiceCollection, error) {
	files, err := os.ReadDir(sourcePath)
	if err != nil {
		return nil, fmt.Errorf("cannot read list of files from %s: %v", sourcePath, err)
	}

	collection := &invoiceCollection{
		Files: make([]string, 0),
	}

	var fileName string

	for _, file := range files {
		fileName = file.Name()

		if strings.HasPrefix(fileName, "invoice-") && filepath.Ext(fileName) == ".xml" {
			collection.Files = append(collection.Files, filepath.Join(sourcePath, fileName))

			if collection.Issuer == "" {
				if collection.Issuer, err = parseInvoiceIssuer(collection.Files[len(collection.Files)-1]); err != nil {
					return nil, fmt.Errorf("cannot parse invoice issuer from %s: %v", fileName, err)
				}
			}
		}
	}

	if collection.Issuer == "" {
		return nil, fmt.Errorf("no issuer found - this should be impossible")
	}

	return collection, nil

}
