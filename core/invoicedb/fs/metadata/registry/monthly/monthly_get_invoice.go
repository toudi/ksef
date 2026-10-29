package monthly

import "ksef/core/invoicedb"

func (r *Registry) GetInvoice(filter func(i *Invoice) bool) (*Invoice, error) {
	for _, invoice := range r.state.Invoices {
		if filter(invoice) {
			return invoice, nil
		}
	}

	return nil, invoicedb.ErrNotFound
}

func (r *Registry) Invoices() []*Invoice {
	return r.state.Invoices
}
