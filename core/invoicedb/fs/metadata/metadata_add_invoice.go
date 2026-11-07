package metadata

import (
	"context"
	"ksef/core/invoicedb"
	"ksef/core/ksef"
)

func (fs *FSInvoiceMetadataStorage) AddReceivedInvoice(
	ctx context.Context,
	yearMonth invoicedb.YearMonth,
	subjectType ksef.SubjectType,
	meta ksef.InvoiceMetadata,
	qrCodeBaseURL string,
) (*invoicedb.StoredInvoice, error) {
	invoiceRegistry, err := fs.getRegistryByYearMonth(yearMonth)
	if err != nil {
		return nil, err
	}
	return invoiceRegistry.AddReceivedInvoice(subjectType, meta, fs.nip, qrCodeBaseURL)
}
