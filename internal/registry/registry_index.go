package registry

import (
	"errors"
	"ksef/internal/registry/types"
)

var (
	ErrUnknownInvoice = errors.New("unable to lookup invoice by checksum")
)

func (r *InvoiceRegistry) Contains(ksefRefNo string) bool {
	_, exists := r.seiRefNoIndex[ksefRefNo]
	return exists
}

func (r *InvoiceRegistry) ContainsHash(hash string) bool {
	_, exists := r.checksumIndex[hash]
	return exists
}

func (r *InvoiceRegistry) GetSEIRefNo(invoiceNo string) (string, error) {
	for _, invoice := range r.Invoices {
		if invoice.ReferenceNumber == invoiceNo || invoice.KSeFReferenceNumber == invoiceNo {
			return invoice.KSeFReferenceNumber, nil
		}
	}

	return "", errors.New("invoice number could not be found")
}

func (r *InvoiceRegistry) GetInvoiceByChecksum(checksum string) (types.Invoice, error) {
	index, exists := r.checksumIndex[checksum]
	if !exists {
		return types.Invoice{}, ErrUnknownInvoice
	}
	return r.Invoices[index], nil
}
