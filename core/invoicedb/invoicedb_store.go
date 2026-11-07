package invoicedb

import (
	"context"
	"io"
	"ksef/core/invoicedb/invoice"
	"ksef/core/ksef"
)

// StoreReceivedInvoice registers received invoice metadata and stores the XML
// payload. The metadata is mutated in memory first so that the storage ID can
// be derived; if the data store fails the caller must treat the whole command
// as failed, which prevents InvoiceDB.Close from persisting the metadata.
func (i *InvoiceDB) StoreReceivedInvoice(
	ctx context.Context,
	yearMonth YearMonth,
	subjectType ksef.SubjectType,
	meta ksef.InvoiceMetadata,
	content io.Reader,
	qrCodeBaseURL string,
) (*StoredInvoice, error) {
	stored, err := i.metadataStorage.AddReceivedInvoice(ctx, yearMonth, subjectType, meta, qrCodeBaseURL)
	if err != nil {
		return nil, err
	}

	if _, err := i.dataStorage.Store(
		ctx,
		content,
		stored.StorageId(invoice.InvoiceFileXML),
		false,
	); err != nil {
		return nil, err
	}

	return stored, nil
}
