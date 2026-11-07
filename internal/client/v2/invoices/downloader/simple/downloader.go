package simple

import (
	"context"
	"ksef/core/downloader"
	"ksef/core/invoicedb"
	"ksef/internal/http"
)

// InvoiceHashChecker is the minimal storage lookup the simple downloader needs.
// Both the new core InvoiceDB and the legacy monthly registry adapter satisfy it.
type InvoiceHashChecker interface {
	ContainsHash(ctx context.Context, yearMonth invoicedb.YearMonth, checksum string) (bool, error)
}

type SimpleDownloader struct {
	httpClient  *http.Client
	hashChecker InvoiceHashChecker
	params      downloader.Config
}

func NewDownloader(
	httpClient *http.Client,
	hashChecker InvoiceHashChecker,
) *SimpleDownloader {
	return &SimpleDownloader{
		httpClient:  httpClient,
		hashChecker: hashChecker,
	}
}

func (d *SimpleDownloader) Close() error {
	return nil
}
