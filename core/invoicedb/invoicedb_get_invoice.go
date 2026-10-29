package invoicedb

import (
	"context"
	"ksef/internal/logging"
	"time"
)

// GetInvoiceByRefNo queries InvoiceDB's metadata storage for the given invoice.
// if the invoice can be found, InvoiceDB returns the corresponding storage ID.
// the storage ID is only relevant for the filesystem implementation as it
// needs to persist the in-memory buffers to disk.
func (i *InvoiceDB) GetInvoiceByRefNo(
	ctx context.Context,
	refNo string,
) (invoice *StoredInvoice, err error) {
	now := time.Now()
	now = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)

	for range i.maxMonthsLookup {
		yearMonth := YearMonth{Year: now.Year(), Month: int(now.Month())}

		logging.InvoicesDBLogger.Debug("trying to find invoice in storage", "yearMonth", yearMonth)
		invoice, err = i.metadataStorage.GetInvoiceByUID(
			ctx,
			InvoiceUID{
				YearMonth: yearMonth,
				RefNo:     &refNo,
			},
		)
		if err != nil && err != ErrNotFound {
			return nil, err
		}
		if invoice != nil && err == nil {
			return invoice, err
		}
		// go back 1 month
		now = now.AddDate(0, -1, 0)
	}
	return nil, ErrNotFound
}
