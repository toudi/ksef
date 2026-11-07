package invoicedb

import (
	"context"
	"io"
	"time"

	"ksef/core/invoicedb/invoice"
	"ksef/core/ksef"
)

type YearMonth struct {
	Year  int
	Month int
}

type InvoiceUID struct {
	NIP       string
	YearMonth YearMonth
	RefNo     *string
	KSeFRefNo *string
}

type StorageId struct {
	NIP       string
	YearMonth YearMonth
	InvoiceID InvoiceID
	FileType  invoice.InvoiceFileType
}

type StoredInvoice struct {
	NIP       string
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
		NIP:      si.NIP,
		FileType: fileType,
	}
}

type InvoiceID struct {
	ID          string
	InvoiceType invoice.InvoiceType
}

type InvoiceSyncParams struct {
	LastTimestamp time.Time
	SubjectTypes  []ksef.SubjectType
}

type InvoiceSyncState struct {
	YearMonth     YearMonth
	LastTimestamp time.Time
}

// InvoiceMetadataStorage is responsible for keeping information about the
// invoice metadata, such as:
// - it's checksum
// - reference numbers (KSeF ref no, invoice reference number)
// - corrections
// - etc.
type InvoiceMetadataStorage interface {
	GetInvoiceByUID(ctx context.Context, invoiceUID InvoiceUID) (invoice *StoredInvoice, err error)

	ContainsHash(ctx context.Context, yearMonth YearMonth, checksum string) (bool, error)

	GetSyncParams(ctx context.Context, yearMonth YearMonth) (*InvoiceSyncParams, error)
	SetSyncParams(ctx context.Context, yearMonth YearMonth, params InvoiceSyncParams) error

	// GetLastSyncState returns the newest month with non-zero sync state,
	// together with its LastTimestamp. nil means no synced month was found.
	GetLastSyncState(ctx context.Context) (*InvoiceSyncState, error)

	AddReceivedInvoice(
		ctx context.Context,
		yearMonth YearMonth,
		subjectType ksef.SubjectType,
		meta ksef.InvoiceMetadata,
		qrCodeBaseURL string,
	) (*StoredInvoice, error)

	Close(ctx context.Context) error
}

// InvoiceStorage is responsible for storing invoice data files (XML files)
type InvoiceDataStorage interface {
	Open(ctx context.Context, storageId StorageId) (io.ReadCloser, error)
	Store(ctx context.Context, src io.Reader, storageId StorageId, force bool) (bool, error)

	// LocalPath returns a filesystem path usable by external tools.
	//
	// ok == true: path is the authoritative storage location; writes to it are
	// final and do not need Store. The parent directory is guaranteed to exist.
	// ok == false: the caller must stage through a workdir and call Store.
	LocalPath(storageId StorageId) (path string, ok bool)

	Close(ctx context.Context) error
}
