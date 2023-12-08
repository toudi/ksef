package invoice

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
	var fullFileName string
	var issuer string

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
		if issuer, _ = parseInvoiceIssuer(fullFileName); issuer != "" {
			collection.Files = append(collection.Files, fullFileName)
			if collection.Issuer == "" {
				collection.Issuer = issuer
			}
		}
	}

	if collection.Issuer == "" {
		return nil, fmt.Errorf("no issuer found - this should be impossible")
	}

	return collection, nil

}
