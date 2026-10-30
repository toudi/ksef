package smart

import (
	"bytes"
	"context"
	"ksef/core/ksef"
	"ksef/internal/client/v2/invoices/downloader/simple"
	"ksef/internal/client/v2/types/invoices"
)

func (sd *smartDownloader) downloadWithSimple(
	ctx context.Context,
	subjectType ksef.SubjectType,
	invoicesMetadata *invoices.InvoiceMetadataResponse,
	invoiceReady func(
		subjectType ksef.SubjectType,
		invoice invoices.InvoiceMetadata,
		content bytes.Buffer,
	) error,
) error {
	if sd.simpleDownloader == nil {
		sd.simpleDownloader = simple.NewDownloader(
			sd.httpClient, sd.registry, sd.params,
		)
	}
	return sd.simpleDownloader.DownloadInvoices(ctx, invoicesMetadata, subjectType, invoiceReady)
}
