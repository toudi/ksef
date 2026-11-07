package v2

import (
	"ksef/core/certsdb"
	"ksef/core/downloader"
	"ksef/internal/client/v2/invoices/downloader/export"
	downloaderinterface "ksef/internal/client/v2/invoices/downloader/interface"
	"ksef/internal/client/v2/invoices/downloader/legacy"
	"ksef/internal/client/v2/invoices/downloader/simple"
	"ksef/internal/client/v2/invoices/downloader/smart"
	monthlyregistry "ksef/internal/invoicesdb/monthly-registry"

	"github.com/spf13/viper"
)

// InvoiceDownloader preserves the legacy registry-backed download path used by
// the not-yet-migrated InvoicesDB flows (sync, parallel download).
func (c *APIClient) InvoiceDownloader(
	vip *viper.Viper,
	certsDB *certsdb.CertsDB,
	params downloader.Config,
	registry *monthlyregistry.Registry,
) downloaderinterface.InvoiceDownloader {
	if params.UseSmartMode {
		return smart.NewDownloader(
			vip, certsDB, c.authenticatedHTTPClient(), registry,
		)
	}
	if params.UseExportMode {
		return export.NewDownloader(
			vip, certsDB, c.authenticatedHTTPClient(), registry,
		)
	}
	return simple.NewDownloader(
		c.authenticatedHTTPClient(),
		legacy.RegistryHashChecker{Registry: registry},
	)
}
