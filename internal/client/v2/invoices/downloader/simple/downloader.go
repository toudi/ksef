package simple

import (
	"ksef/internal/client/v2/types/invoices"
	"ksef/internal/http"
	monthlyregistry "ksef/internal/invoicesdb/monthly-registry"
)

type SimpleDownloader struct {
	httpClient *http.Client
	registry   *monthlyregistry.Registry
	params     invoices.DownloadParams
}

func NewDownloader(
	httpClient *http.Client,
	registry *monthlyregistry.Registry,
	params invoices.DownloadParams,
) *SimpleDownloader {
	return &SimpleDownloader{
		httpClient: httpClient,
		registry:   registry,
		params:     params,
	}
}

func (d *SimpleDownloader) Close() error {
	return nil
}
