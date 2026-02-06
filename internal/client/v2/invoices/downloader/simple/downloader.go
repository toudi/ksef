package simple

import (
	downloaderinterface "ksef/internal/client/v2/invoices/downloader/interface"
	"ksef/internal/client/v2/types/invoices"
	"ksef/internal/http"
	monthlyregistry "ksef/internal/invoicesdb/monthly-registry"
)

type simpleDownloader struct {
	httpClient *http.Client
	registry   *monthlyregistry.Registry
	params     invoices.DownloadParams
}

func NewDownloader(
	httpClient *http.Client,
	registry *monthlyregistry.Registry,
	params invoices.DownloadParams,
) downloaderinterface.InvoiceDownloader {
	return &simpleDownloader{
		httpClient: httpClient,
		registry:   registry,
		params:     params,
	}
}
