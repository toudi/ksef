package v2

import (
	"ksef/internal/client/v2/invoices"
	invoiceTypes "ksef/internal/client/v2/types/invoices"
	monthlyregistry "ksef/internal/invoicesdb/monthly-registry"
)

func (c *APIClient) InvoiceDownloader(
	params invoiceTypes.DownloadParams,
	registry *monthlyregistry.Registry,
) *invoices.InvoiceDownloader {
	return invoices.NewInvoiceDownloader(
		c.authenticatedHTTPClient(),
		params,
		registry,
	)
}
