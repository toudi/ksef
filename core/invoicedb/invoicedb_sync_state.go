package invoicedb

import "context"

func (i *InvoiceDB) GetSyncParams(
	ctx context.Context,
	yearMonth YearMonth,
) (*InvoiceSyncParams, error) {
	return i.metadataStorage.GetSyncParams(ctx, yearMonth)
}

func (i *InvoiceDB) SetSyncParams(
	ctx context.Context,
	yearMonth YearMonth,
	params InvoiceSyncParams,
) error {
	return i.metadataStorage.SetSyncParams(ctx, yearMonth, params)
}

func (i *InvoiceDB) GetLastSyncState(ctx context.Context) (*InvoiceSyncState, error) {
	return i.metadataStorage.GetLastSyncState(ctx)
}
