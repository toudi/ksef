package v2

import (
	"ksef/core/certsdb"
	"ksef/core/downloader"
	"ksef/internal/client/v2/invoices"
	downloaderinterface "ksef/internal/client/v2/invoices/downloader/interface"
	monthlyregistry "ksef/internal/invoicesdb/monthly-registry"

	"github.com/spf13/viper"
)

func (c *APIClient) InvoiceDownloader(
	vip *viper.Viper,
	certsDB *certsdb.CertsDB,
	params downloader.Config,
	registry *monthlyregistry.Registry,
) downloaderinterface.InvoiceDownloader {
	return invoices.NewInvoiceDownloader(
		vip,
		certsDB,
		c.authenticatedHTTPClient(),
		params,
		registry,
	)
}
