package v2

import (
	"ksef/internal/certsdb"
	"ksef/internal/client/v2/invoices"
	downloaderinterface "ksef/internal/client/v2/invoices/downloader/interface"
	invoiceTypes "ksef/internal/client/v2/types/invoices"
	monthlyregistry "ksef/internal/invoicesdb/monthly-registry"

	"github.com/spf13/viper"
)

func (c *APIClient) InvoiceDownloader(
	vip *viper.Viper,
	certsDB *certsdb.CertificatesDB,
	params invoiceTypes.DownloadParams,
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
