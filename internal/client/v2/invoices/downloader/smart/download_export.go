package smart

import (
	"bytes"
	"context"
	"ksef/core/ksef"
	"ksef/internal/client/v2/invoices/downloader/export"
)

func (sd *smartDownloader) downloadWithExport(
	ctx context.Context,
	invoiceReady func(
		subjectType ksef.SubjectType,
		invoice ksef.InvoiceMetadata,
		content bytes.Buffer,
	) error,
) error {
	if sd.exportDownloader == nil {
		sd.exportDownloader = export.NewDownloader(
			sd.vip,
			sd.certsDB,
			sd.httpClient,
			sd.registry,
		)
	}
	return sd.exportDownloader.Download(ctx, sd.params, invoiceReady)
}
