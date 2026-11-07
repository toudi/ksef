package registry

import (
	"errors"
	"ksef/core/invoicedb"
	"ksef/core/invoicedb/fs/metadata/registry/annual"
	"ksef/core/invoicedb/fs/metadata/registry/monthly"
	"ksef/core/ksef"
	"ksef/internal/logging"
)

type FSInvoiceRegistry struct {
	root      string
	yearMonth invoicedb.YearMonth
	nip       string

	annualRegistry  *annual.Registry
	monthlyRegistry *monthly.Registry
}

// NIP is kept here because when we eventually add the invoice to registry, we need to return back the NIP.
// obviously we could parse the NIP from `root` but why bother
func Open(root string, yearMonth invoicedb.YearMonth, nip string) (*FSInvoiceRegistry, error) {
	logging.InvoicesDBLogger.Debug("fs/metadata/registry::Open", "root", root, "yearMonth", yearMonth)
	annualRegistry, err := annual.Open(root, yearMonth)
	if err != nil {
		return nil, err
	}

	monthlyRegistry, err := monthly.Open(root, yearMonth)
	if err != nil {
		return nil, err
	}

	return &FSInvoiceRegistry{
		root:            root,
		yearMonth:       yearMonth,
		nip:             nip,
		annualRegistry:  annualRegistry,
		monthlyRegistry: monthlyRegistry,
	}, nil
}

func (fs *FSInvoiceRegistry) Close() error {
	return errors.Join(
		fs.monthlyRegistry.Close(),
		fs.annualRegistry.Close(),
	)
}

func (fs *FSInvoiceRegistry) ContainsHash(checksum string) bool {
	return fs.monthlyRegistry.ContainsHash(checksum)
}

func (fs *FSInvoiceRegistry) GetSyncParams() *invoicedb.InvoiceSyncParams {
	return fs.monthlyRegistry.SyncParams()
}

func (fs *FSInvoiceRegistry) SetSyncParams(params invoicedb.InvoiceSyncParams) {
	fs.monthlyRegistry.SetSyncParams(params)
}

func (fs *FSInvoiceRegistry) AddReceivedInvoice(
	subjectType ksef.SubjectType,
	meta ksef.InvoiceMetadata,
	nip string,
	qrCodeBaseURL string,
) (*invoicedb.StoredInvoice, error) {
	return fs.monthlyRegistry.AddReceivedInvoice(
		subjectType,
		meta,
		nip,
		qrCodeBaseURL,
	)
}
