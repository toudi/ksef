package uploader

import "ksef/internal/invoice"

func (idb *InvoiceDB) Add(i *invoice.Invoice, checksum string) {
	idb.Invoices = append(idb.Invoices, Invoice{
		RefNo:          i.Number,
		Checksum:       checksum,
		GenerationTime: i.GenerationTime,
	})
}
