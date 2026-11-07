package invoicedb

import "context"

func (i *InvoiceDB) ContainsHash(
	ctx context.Context,
	yearMonth YearMonth,
	checksum string,
) (bool, error) {
	return i.metadataStorage.ContainsHash(ctx, yearMonth, checksum)
}
