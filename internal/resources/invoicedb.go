package resources

import (
	"ksef/core/invoicedb"
	fsInvoiceData "ksef/core/invoicedb/fs/data"
	fsInvoiceMeta "ksef/core/invoicedb/fs/metadata"
	"ksef/internal/invoicesdb/config"
	"ksef/internal/runtime"
	"path/filepath"

	"github.com/spf13/viper"
)

var (
	invoiceDBInstance                *invoicedb.InvoiceDB
	invoiceFSDataStorageInstance     *fsInvoiceData.FSDataStorage
	invoiceFSMetadataStorageInstance *fsInvoiceMeta.FSInvoiceMetadataStorage
)

func InvoiceDataStorage(vip *viper.Viper) (*fsInvoiceData.FSDataStorage, error) {
	if invoiceFSDataStorageInstance != nil {
		return invoiceFSDataStorageInstance, nil
	}

	root, err := InvoiceStorageRoot(vip)
	if err != nil {
		return nil, err
	}

	invoiceFSDataStorageInstance = fsInvoiceData.NewFSDataStorage(root)
	return invoiceFSDataStorageInstance, nil
}

func InvoiceMetadataStorage(vip *viper.Viper) (*fsInvoiceMeta.FSInvoiceMetadataStorage, error) {
	if invoiceFSMetadataStorageInstance != nil {
		return invoiceFSMetadataStorageInstance, nil
	}

	root, err := InvoiceStorageRoot(vip)
	if err != nil {
		return nil, err
	}

	invoiceFSMetadataStorageInstance = fsInvoiceMeta.NewFSInvoiceMetadataStorage(root)
	return invoiceFSMetadataStorageInstance, nil
}

func InvoiceStorageRoot(vip *viper.Viper) (string, error) {
	cfg := config.GetInvoicesDBConfig(vip)
	environmentId := runtime.GetEnvironmentId(vip)
	nip, err := runtime.GetNIP(vip)
	if err != nil {
		return "", err
	}

	return filepath.Join(
		cfg.Root,
		environmentId,
		nip,
	), nil
}

func InvoiceDB(vip *viper.Viper) (*invoicedb.InvoiceDB, error) {
	if invoiceDBInstance != nil {
		return invoiceDBInstance, nil
	}

	invoiceMetadataStorage, err := InvoiceMetadataStorage(vip)
	if err != nil {
		return nil, err
	}
	invoiceDataStorage, err := InvoiceDataStorage(vip)
	if err != nil {
		return nil, err
	}

	invoiceDB, err := invoicedb.NewInvoiceDB(invoiceMetadataStorage, invoiceDataStorage)
	if err != nil {
		return nil, err
	}

	invoiceDBInstance = invoiceDB
	return invoiceDBInstance, nil
}
