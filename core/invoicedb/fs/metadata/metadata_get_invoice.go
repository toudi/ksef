package metadata

import (
	"context"
	"errors"
	"ksef/core/invoicedb"
)

func (fs *FSInvoiceMetadataStorage) GetInvoiceByUID(
	ctx context.Context,
	invoiceUID invoicedb.InvoiceUID,
) (*invoicedb.StoredInvoice, error) {
	// retrieve metadata database by storage id
	invoiceRegistry, err := fs.getRegistryByYearMonth(invoiceUID.YearMonth)
	if err != nil {
		return nil, errors.Join(err, invoicedb.ErrNotFound)
	}

	return invoiceRegistry.GetInvoiceByUID(invoiceUID)
}
