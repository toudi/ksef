package invoicedb

import (
	"context"
	"io"
	"ksef/core/invoicedb/invoice"
)

type YearMonth struct {
	Year  int
	Month int
}

type InvoiceUID struct {
	YearMonth YearMonth
	RefNo     *string
	KSeFRefNo *string
}

type StorageId struct {
	YearMonth YearMonth
	InvoiceID InvoiceID
	FileType  invoice.InvoiceFileType
}

type StoredInvoice struct {
	ID        string
	Invoice   *invoice.Invoice
	YearMonth YearMonth
}

func (si StoredInvoice) StorageId(fileType invoice.InvoiceFileType) StorageId {
	return StorageId{
		YearMonth: si.YearMonth,
		InvoiceID: InvoiceID{
			ID:          si.ID,
			InvoiceType: si.Invoice.Type,
		},
		FileType: fileType,
	}
}

type InvoiceID struct {
	ID          string
	InvoiceType invoice.InvoiceType
}

// InvoiceMetadataStorage is responsible for keeping information about the
// invoice metadata, such as:
// - it's checksum
// - reference numbers (KSeF ref no, invoice reference number)
// - corrections
// - etc.
type InvoiceMetadataStorage interface {
	GetInvoiceByUID(ctx context.Context, invoiceUID InvoiceUID) (invoice *StoredInvoice, err error)
	Close(ctx context.Context) error
}

// InvoiceStorage is responsible for storing invoice data files (XML files)
type InvoiceDataStorage interface {
	Open(ctx context.Context, storageId StorageId) (io.ReadCloser, error)
	Store(ctx context.Context, src io.Reader, storageId StorageId, force bool) (bool, error)
	Close(ctx context.Context) error
}
