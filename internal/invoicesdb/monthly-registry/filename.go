package monthlyregistry

import (
	"fmt"
	"ksef/internal/sei"

	"github.com/mozillazg/go-slugify"
)

var dirnameByType = map[InvoiceType]string{
	InvoiceTypeIssued:   "wystawione",
	InvoiceTypeReceived: "otrzymane",
}

func (r *Registry) getIssuedInvoiceFilename(invoiceNumber string, ordNo int) string {
	return fmt.Sprintf(
		"%s/%s/%04d-%s.xml",
		r.dir,
		dirnameByType[InvoiceTypeIssued],
		ordNo,
		slugify.Slugify(invoiceNumber),
	)
}

func (r *Registry) GetDestFileName(inv *sei.ParsedInvoice, invoiceType InvoiceType) string {
	numInvoices := r.countInvoicesByType(invoiceType)

	if invoiceType == InvoiceTypeIssued {
		return r.getIssuedInvoiceFilename(inv.Invoice.Number, numInvoices+1)
	}

	return fmt.Sprintf(
		"%s/%s/%04d-%s-%s.xml",
		r.dir,
		dirnameByType[invoiceType],
		numInvoices+1,
		slugify.Slugify(inv.Invoice.Issuer.Name),
		slugify.Slugify(inv.Invoice.Number),
	)

}

func (r *Registry) countInvoicesByType(invoiceType InvoiceType) (count int) {
	for _, invoice := range r.invoices {
		if invoice.Type == invoiceType {
			count += 1
		}
	}

	return count
}
