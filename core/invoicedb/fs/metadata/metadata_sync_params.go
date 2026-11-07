package metadata

import (
	"context"
	"time"

	"ksef/core/invoicedb"
)

func (fs *FSInvoiceMetadataStorage) GetSyncParams(
	ctx context.Context,
	yearMonth invoicedb.YearMonth,
) (*invoicedb.InvoiceSyncParams, error) {
	invoiceRegistry, err := fs.getRegistryByYearMonth(yearMonth)
	if err != nil {
		return nil, err
	}
	return invoiceRegistry.GetSyncParams(), nil
}

func (fs *FSInvoiceMetadataStorage) SetSyncParams(
	ctx context.Context,
	yearMonth invoicedb.YearMonth,
	params invoicedb.InvoiceSyncParams,
) error {
	invoiceRegistry, err := fs.getRegistryByYearMonth(yearMonth)
	if err != nil {
		return err
	}
	invoiceRegistry.SetSyncParams(params)
	return nil
}

func (fs *FSInvoiceMetadataStorage) GetLastSyncState(
	ctx context.Context,
) (*invoicedb.InvoiceSyncState, error) {
	now := time.Now().Local()
	month := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	for range 12 {
		yearMonth := invoicedb.YearMonth{Year: month.Year(), Month: int(month.Month())}
		invoiceRegistry, err := fs.getRegistryByYearMonth(yearMonth)
		if err != nil {
			return nil, err
		}
		params := invoiceRegistry.GetSyncParams()
		if params != nil && !params.LastTimestamp.IsZero() {
			return &invoicedb.InvoiceSyncState{
				YearMonth:     yearMonth,
				LastTimestamp: params.LastTimestamp,
			}, nil
		}
		month = month.AddDate(0, -1, 0)
	}

	return nil, nil
}
