package metadata

import (
	"context"
	"errors"

	"ksef/core/invoicedb"
	"ksef/core/invoicedb/fs/metadata/registry"
	"ksef/internal/logging"
)

var _ invoicedb.InvoiceMetadataStorage = (*FSInvoiceMetadataStorage)(nil)

type FSInvoiceMetadataStorage struct {
	root      string
	nip       string

	invoiceRegistries map[invoicedb.YearMonth]*registry.FSInvoiceRegistry
}

func NewFSInvoiceMetadataStorage(
	root string,
	nip string,
) *FSInvoiceMetadataStorage {
	logging.InvoicesDBLogger.Debug("metadata storage", "root", root, "nip", nip)
	return &FSInvoiceMetadataStorage{
		root:              root,
		nip:               nip,
		invoiceRegistries: make(map[invoicedb.YearMonth]*registry.FSInvoiceRegistry),
	}
}

func (fs *FSInvoiceMetadataStorage) Close(ctx context.Context) error {
	var err error
	for _, invoiceRegistry := range fs.invoiceRegistries {
		err = errors.Join(err, invoiceRegistry.Close())
	}
	return err
}
