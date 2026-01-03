package monthlyregistry

import "errors"

var errUnableToFindInvoice = errors.New("unable to lookup invoice")

func (r *Registry) UpdateInvoiceByChecksum(checksum string, modify func(invoice *Invoice) error) error {
	invoice := r.GetInvoiceByChecksum(checksum)
	if invoice == nil {
		return errUnableToFindInvoice
	}

	return modify(invoice)
}
