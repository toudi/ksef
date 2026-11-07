package metadata

import (
	"ksef/core/invoicedb"
	"ksef/core/invoicedb/fs/metadata/registry"
)

func (fs *FSInvoiceMetadataStorage) getRegistryByYearMonth(
	yearMonth invoicedb.YearMonth,
) (*registry.FSInvoiceRegistry, error) {
	// check if it is already opened
	if registry, exists := fs.invoiceRegistries[yearMonth]; exists {
		return registry, nil
	}
	// fallback to opening the registry
	registry, err := registry.Open(fs.root, yearMonth, fs.nip)
	if err != nil {
		return nil, err
	}
	// we're fine - just store the registry in cache
	fs.invoiceRegistries[yearMonth] = registry
	return registry, nil
}
