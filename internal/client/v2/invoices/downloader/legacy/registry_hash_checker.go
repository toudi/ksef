package legacy

import (
	"context"

	"ksef/core/invoicedb"
	monthlyregistry "ksef/internal/invoicesdb/monthly-registry"
)

// RegistryHashChecker adapts the legacy monthly registry to the simple
// downloader's hash-checking interface.
type RegistryHashChecker struct {
	Registry *monthlyregistry.Registry
}

func (c RegistryHashChecker) ContainsHash(
	_ context.Context,
	_ invoicedb.YearMonth,
	checksum string,
) (bool, error) {
	return c.Registry.ContainsHash(checksum), nil
}
