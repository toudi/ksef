package registry

import "ksef/core/invoicedb/fs/metadata/registry/monthly"

func (fs *FSInvoiceRegistry) Invoices() []*monthly.Invoice {
	return fs.monthlyRegistry.Invoices()
}
