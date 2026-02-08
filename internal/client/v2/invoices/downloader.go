package invoices

import (
	"ksef/internal/certsdb"
	"ksef/internal/client/v2/invoices/downloader/export"
	downloaderinterface "ksef/internal/client/v2/invoices/downloader/interface"
	"ksef/internal/client/v2/invoices/downloader/simple"
	"ksef/internal/client/v2/types/invoices"
	"ksef/internal/http"
	monthlyregistry "ksef/internal/invoicesdb/monthly-registry"

	"github.com/spf13/viper"
)

func NewInvoiceDownloader(
	vip *viper.Viper,
	certsDB *certsdb.CertificatesDB,
	httpClient *http.Client,
	downloadParams invoices.DownloadParams,
	registry *monthlyregistry.Registry,
) downloaderinterface.InvoiceDownloader {
	if downloadParams.UseExportMode {
		return export.NewDownloader(
			vip, certsDB, httpClient, registry, downloadParams,
		)
	}
	return simple.NewDownloader(
		httpClient, registry, downloadParams,
	)
}
