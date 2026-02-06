package v2

import (
	"ksef/internal/certsdb"
	"ksef/internal/client/v2/invoices"
	downloaderinterface "ksef/internal/client/v2/invoices/downloader/interface"
	invoiceTypes "ksef/internal/client/v2/types/invoices"
	monthlyregistry "ksef/internal/invoicesdb/monthly-registry"
)

func (c *APIClient) InvoiceDownloader(
	certsDB *certsdb.CertificatesDB,
	params invoiceTypes.DownloadParams,
	registry *monthlyregistry.Registry,
) downloaderinterface.InvoiceDownloader {
	return invoices.NewInvoiceDownloader(
		certsDB,
		c.authenticatedHTTPClient(),
		params,
		registry,
	)
}
