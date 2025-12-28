package monthlyregistry

import (
	"errors"
	"fmt"
	"io/fs"
	"ksef/internal/logging"
	"path/filepath"
	"strings"

	"github.com/mozillazg/go-slugify"
)

func (r *Registry) postOpenHook() {
	var assignOrdNums bool = len(r.OrdNums) == 0

	for index, invoice := range r.Invoices {
		invoiceType := invoice.Type

		if assignOrdNums {
			if _, exists := r.OrdNums[invoiceType]; !exists {
				r.OrdNums[invoiceType] = 0
			}
			r.OrdNums[invoiceType]++
			invoice.OrdNum = r.OrdNums[invoiceType]
		}

		if invoiceType > InvoiceTypeIssued && invoice.Issuer == nil {
			guessedFilename, err := r.getInvoiceFilenameHeuristic(invoice)
			if err != nil {
				logging.InvoicesDBLogger.Error("unable to guess invoice filename", "err", err)
			} else {
				xmlInvoice, _, err := ParseInvoice(guessedFilename)
				if err != nil {
					logging.InvoicesDBLogger.Error("unable to parse XML invoice", "err", err)
					continue
				}
				invoice.Issuer = &InvoiceIssuer{
					NIP:  xmlInvoice.Issuer,
					Name: xmlInvoice.IssuerName,
				}
			}
		}

		r.checksumIndex[invoice.Checksum] = index
	}
}

var errSkipTraversal = errors.New("skip traversal")

func (r *Registry) getInvoiceFilenameHeuristic(i *Invoice) (guessedFilename string, err error) {
	dirName := dirnameByType[i.Type]

	if err = filepath.WalkDir(filepath.Join(r.dir, dirName), func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if filepath.Ext(d.Name()) != ".xml" {
			return nil
		}

		baseName := filepath.Base(d.Name())

		if strings.HasPrefix(baseName, fmt.Sprintf("%04d", i.OrdNum)) && strings.Contains(baseName, slugify.Slugify(i.RefNo)) {
			guessedFilename = path
			return errSkipTraversal
		}

		return nil
	}); err != nil {
		if err == errSkipTraversal {
			return guessedFilename, nil
		}
	}

	return guessedFilename, nil
}
