package metadata

import (
	"ksef/core/invoicedb"
	"ksef/core/invoicedb/fs/metadata/registry"
	"path/filepath"
	"strconv"
)

func (fs *FSInvoiceMetadataStorage) getRegistryByYearMonth(
	yearMonth invoicedb.YearMonth,
) (*registry.FSInvoiceRegistry, error) {
	// check if it is already opened
	if registry, exists := fs.invoiceRegistries[yearMonth]; exists {
		return registry, nil
	}
	// fallback to opening the registry
	// calculate the registry path based on root
	registryRoot := fs.getRegistryRoot(yearMonth)
	registry, err := registry.Open(registryRoot, yearMonth)
	if err != nil {
		return nil, err
	}
	// we're fine - just store the registry in cache
	fs.invoiceRegistries[yearMonth] = registry
	return registry, nil
}

func (fs *FSInvoiceMetadataStorage) getRegistryRoot(yearMonth invoicedb.YearMonth) string {
	return filepath.Join(
		fs.root,
		strconv.Itoa(yearMonth.Year),
	)
}
