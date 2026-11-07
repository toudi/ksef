package invoices

import (
	"ksef/core/certsdb"
	"ksef/core/downloader"
	"ksef/core/invoicedb"
	downloaderinterface "ksef/internal/client/v2/invoices/downloader/interface"
	"ksef/internal/client/v2/invoices/downloader/simple"
	"ksef/internal/http"
)

func NewInvoiceDownloader(
	config downloader.Config,
	certsDB *certsdb.CertsDB,
	httpClient *http.Client,
	invoiceDB *invoicedb.InvoiceDB,
) downloaderinterface.InvoiceDownloader {
	// if downloadParams.UseSmartMode {
	// 	return smart.NewDownloader(
	// 		vip, certsDB, httpClient, registry, downloadParams,
	// 	)
	// }
	// if downloadParams.UseExportMode {
	// 	return export.NewDownloader(
	// 		vip, certsDB, httpClient, registry, downloadParams,
	// 	)
	// }
	return simple.NewDownloader(
		httpClient, invoiceDB,
	)
}
