package export

import (
	"ksef/internal/certsdb"
	downloaderinterface "ksef/internal/client/v2/invoices/downloader/interface"
	"ksef/internal/client/v2/types/invoices"
	"ksef/internal/http"
	monthlyregistry "ksef/internal/invoicesdb/monthly-registry"
)

type exportDownloader struct {
	httpClient *http.Client
	registry   *monthlyregistry.Registry
	params     invoices.DownloadParams
	certsDB    *certsdb.CertificatesDB
}

func NewDownloader(
	httpClient *http.Client,
	registry *monthlyregistry.Registry,
	params invoices.DownloadParams,
) downloaderinterface.InvoiceDownloader {
	return &exportDownloader{
		httpClient: httpClient,
		registry:   registry,
		params:     params,
	}
}
