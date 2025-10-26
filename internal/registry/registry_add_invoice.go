package registry

import "ksef/internal/sei/api/client/v2/types/invoices"

func (r *InvoiceRegistry) AddInvoice(invoice invoices.InvoiceMetadata, checksum string) {
	r.Invoices = append(r.Invoices, Invoice{
		ReferenceNumber:     invoice.InvoiceNumber,
		KSeFReferenceNumber: invoice.KSeFNumber,
		InvoiceType:         invoice.InvoiceType,
		Checksum:            checksum,
	})
	r.seiRefNoIndex[invoice.KSeFNumber] = len(r.Invoices)
	r.checksumIndex[checksum] = len(r.Invoices)
}
