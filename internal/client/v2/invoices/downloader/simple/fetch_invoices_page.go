package simple

import (
	"bytes"
	"context"
	"errors"
	"ksef/internal/client/v2/invoices/downloader/metadata"
	"ksef/internal/client/v2/types/invoices"
	types "ksef/internal/client/v2/types/invoices"
	"log/slog"
)

var (
	errUnableToDownloadInvoice     = errors.New("unable to download invoice")
	errProcessingDownloadedInvoice = errors.New("error processing downloaded invoice")
)

func (sd *SimpleDownloader) DownloadInvoices(
	ctx context.Context,
	metadataPage *types.InvoiceMetadataResponse,
	subjectType invoices.SubjectType,
	invoiceReady func(subjectType invoices.SubjectType, invoiceMeta invoices.InvoiceMetadata, content bytes.Buffer) error,
	logger *slog.Logger,
) (err error) {
	var (
		finished             bool
		page                 int
		invoiceContentBuffer bytes.Buffer
	)

	for !finished {
		// fetch all of the invoices from the page.
		for _, invoice := range metadataPage.Invoices {
			if sd.registry.ContainsHash(invoice.Checksum()) {
				logger.Info("Ta faktura już została pobrana", "hash", invoice.Checksum())
				continue
			}

			invoiceContentBuffer.Reset()

			if err = sd.downloadInvoice(ctx, invoice, &invoiceContentBuffer); err != nil {
				return errors.Join(errUnableToDownloadInvoice, err)
			}

			if err = invoiceReady(subjectType, invoice, invoiceContentBuffer); err != nil {
				return errors.Join(errProcessingDownloadedInvoice, err)
			}
		}

		finished = !metadataPage.HasMore

		if metadataPage.HasMore {
			// try to obtain the next page.
			metadataPage, err = metadata.InvoicesMetadataPage(
				ctx,
				sd.httpClient,
				subjectType,
				sd.params,
				page+1,
			)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
