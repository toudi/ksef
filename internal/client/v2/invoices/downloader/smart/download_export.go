package smart

import (
	"bytes"
	"context"
	"ksef/internal/client/v2/invoices/downloader/export"
	"ksef/internal/client/v2/types/invoices"
)

func (sd *smartDownloader) downloadWithExport(
	ctx context.Context,
	invoiceReady func(
		subjectType invoices.SubjectType,
		invoice invoices.InvoiceMetadata,
		content bytes.Buffer,
	) error,
) error {
	if sd.exportDownloader == nil {
		sd.exportDownloader = export.NewDownloader(
			sd.vip,
			sd.certsDB,
			sd.httpClient,
			sd.registry,
			sd.params,
			sd.logger,
		)
	}
	return sd.exportDownloader.Download(ctx, invoiceReady)
}
