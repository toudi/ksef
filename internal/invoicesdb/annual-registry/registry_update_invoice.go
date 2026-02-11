package annualregistry

import "errors"

var errUnableToFindInvoice = errors.New("unable to find invoice by checksum")

func (r *Registry) UpdateInvoiceByChecksum(checksum string, handler func(invoice *Invoice) error) error {
	if invoice := r.GetByChecksum(checksum); invoice != nil {
		return handler(invoice)
	}
	return errUnableToFindInvoice
}
