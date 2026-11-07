package resources

import (
	"context"
	"path/filepath"

	"ksef/core/invoicedb"
	fsInvoiceData "ksef/core/invoicedb/fs/data"
	fsInvoiceMeta "ksef/core/invoicedb/fs/metadata"
	"ksef/internal/invoicesdb/config"
	"ksef/internal/runtime"

	"github.com/spf13/viper"
)

var (
	invoiceDBInstance                map[string]*invoicedb.InvoiceDB
	invoiceFSDataStorageInstance     *fsInvoiceData.FSDataStorage
	invoiceFSMetadataStorageInstance map[string]*fsInvoiceMeta.FSInvoiceMetadataStorage
)

func init() {
	invoiceDBInstance = make(map[string]*invoicedb.InvoiceDB)
	invoiceFSMetadataStorageInstance = make(map[string]*fsInvoiceMeta.FSInvoiceMetadataStorage)
}

func InvoiceDataStorage(vip *viper.Viper) *fsInvoiceData.FSDataStorage {
	mu.Lock()
	defer mu.Unlock()

	return invoiceDataStorageLocked(vip)
}

func invoiceDataStorageLocked(vip *viper.Viper) *fsInvoiceData.FSDataStorage {
	if invoiceFSDataStorageInstance != nil {
		return invoiceFSDataStorageInstance
	}

	root := InvoiceStorageRoot(vip)

	invoiceFSDataStorageInstance = fsInvoiceData.NewFSDataStorage(root)
	return invoiceFSDataStorageInstance
}

func InvoiceMetadataStorage(vip *viper.Viper, NIP string) *fsInvoiceMeta.FSInvoiceMetadataStorage {
	mu.Lock()
	defer mu.Unlock()

	return invoiceMetadataStorageLocked(vip, NIP)
}

func invoiceMetadataStorageLocked(vip *viper.Viper, NIP string) *fsInvoiceMeta.FSInvoiceMetadataStorage {
	if instance, exists := invoiceFSMetadataStorageInstance[NIP]; exists {
		return instance
	}

	root := filepath.Join(InvoiceStorageRoot(vip), NIP)

	instance := fsInvoiceMeta.NewFSInvoiceMetadataStorage(root, NIP)
	invoiceFSMetadataStorageInstance[NIP] = instance
	return instance
}

func InvoiceStorageRoot(vip *viper.Viper) string {
	cfg := config.GetInvoicesDBConfig(vip)
	environmentId := runtime.GetEnvironmentId(vip)

	return filepath.Join(
		cfg.Root,
		environmentId,
	)
}

func InvoiceDB(vip *viper.Viper, NIP string) (*invoicedb.InvoiceDB, error) {
	mu.Lock()
	defer mu.Unlock()

	return invoiceDBLocked(vip, NIP)
}

func invoiceDBLocked(vip *viper.Viper, NIP string) (*invoicedb.InvoiceDB, error) {
	if instance, exists := invoiceDBInstance[NIP]; exists {
		return instance, nil
	}

	invoiceMetadataStorage := invoiceMetadataStorageLocked(vip, NIP)
	invoiceDataStorage := invoiceDataStorageLocked(vip)

	invoiceDB, err := invoicedb.NewInvoiceDB(invoiceMetadataStorage, invoiceDataStorage)
	if err != nil {
		return nil, err
	}

	destructors = append(destructors, func(ctx context.Context, commandErr error) error {
		if commandErr != nil {
			return nil
		}
		return invoiceDB.Close(ctx)
	})

	invoiceDBInstance[NIP] = invoiceDB
	return invoiceDB, nil
}
