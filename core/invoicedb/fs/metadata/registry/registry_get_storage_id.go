package registry

import (
	"context"
	"ksef/core/invoicedb"
	"ksef/core/invoicedb/invoice"
)

func (fs *FSInvoiceRegistry) GetInvoiceStorageId(
	ctx context.Context,
	invoiceUID invoicedb.InvoiceUID,
	fileType invoice.InvoiceFileType,
) (string, error) {
	return "", nil
}
