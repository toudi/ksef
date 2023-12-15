package registry

import "errors"

func (r *InvoiceRegistry) Contains(refNo string) bool {
	_, exists := r.seiRefNoIndex[refNo]
	return exists
}

func (r *InvoiceRegistry) GetSEIRefNo(invoiceNo string) (string, error) {
	for _, invoice := range r.Invoices {
		if invoice.ReferenceNumber == invoiceNo || invoice.SEIReferenceNumber == invoiceNo {
			return invoice.SEIReferenceNumber, nil
		}
	}

	return "", errors.New("invoice number could not be found")
}
