package metadata

import (
	"context"

	"ksef/core/invoicedb"
)

func (fs *FSInvoiceMetadataStorage) ContainsHash(
	ctx context.Context,
	yearMonth invoicedb.YearMonth,
	checksum string,
) (bool, error) {
	invoiceRegistry, err := fs.getRegistryByYearMonth(yearMonth)
	if err != nil {
		return false, err
	}
	return invoiceRegistry.ContainsHash(checksum), nil
}
