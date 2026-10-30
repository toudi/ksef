package simple

import (
	"ksef/core/downloader"
	"ksef/internal/http"
	monthlyregistry "ksef/internal/invoicesdb/monthly-registry"
)

type SimpleDownloader struct {
	httpClient *http.Client
	registry   *monthlyregistry.Registry
	params     downloader.Config
}

func NewDownloader(
	httpClient *http.Client,
	registry *monthlyregistry.Registry,
	params downloader.Config,
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
