package simple

import (
	"ksef/internal/client/v2/types/invoices"
	"ksef/internal/http"
	monthlyregistry "ksef/internal/invoicesdb/monthly-registry"
	"log/slog"
)

type SimpleDownloader struct {
	httpClient *http.Client
	registry   *monthlyregistry.Registry
	params     invoices.DownloadParams
	logger     *slog.Logger
}

func NewDownloader(
	httpClient *http.Client,
	registry *monthlyregistry.Registry,
	params invoices.DownloadParams,
	logger *slog.Logger,
) *SimpleDownloader {
	return &SimpleDownloader{
		httpClient: httpClient,
		registry:   registry,
		params:     params,
		logger:     logger,
	}
}

func (d *SimpleDownloader) Close() error {
	return nil
}
