package invoicedb

import (
	"context"
	"errors"
)

type InvoiceDB struct {
	dataStorage     InvoiceDataStorage
	metadataStorage InvoiceMetadataStorage
	maxMonthsLookup uint8
}

func NewInvoiceDB(
	metadataStorage InvoiceMetadataStorage,
	dataStorage InvoiceDataStorage,
) (*InvoiceDB, error) {
	return &InvoiceDB{
		dataStorage:     dataStorage,
		metadataStorage: metadataStorage,
		maxMonthsLookup: 12,
	}, nil
}

func (i *InvoiceDB) Close(ctx context.Context) error {
	return errors.Join(
		i.dataStorage.Close(ctx),
		i.metadataStorage.Close(ctx),
	)
}
