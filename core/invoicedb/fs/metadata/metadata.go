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
	root string

	invoiceRegistries map[invoicedb.YearMonth]*registry.FSInvoiceRegistry
}

func NewFSInvoiceMetadataStorage(
	root string,
) *FSInvoiceMetadataStorage {
	logging.InvoicesDBLogger.Debug("metadata storage", "root", root)
	return &FSInvoiceMetadataStorage{
		root:              root,
		invoiceRegistries: make(map[invoicedb.YearMonth]*registry.FSInvoiceRegistry),
	}
}

func (fs *FSInvoiceMetadataStorage) Close(ctx context.Context) error {
	var err error
	for _, registry := range fs.invoiceRegistries {
		err = errors.Join(err, registry.Close())
	}
	return err
}
