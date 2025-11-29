package uploader

import (
	"ksef/internal/invoice"
)

func (idb *InvoiceDB) Add(i *invoice.Invoice, checksum string) error {
	invoiceContent, err := marshallInvoice(i)
	if err != nil {
		return err
	}

	idb.Invoices = append(idb.Invoices, &Invoice{
		RefNo:          i.Number,
		Checksum:       checksum,
		GenerationTime: i.GenerationTime,
		Contents:       invoiceContent,
	})

	return nil
}
