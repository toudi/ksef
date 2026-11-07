package simple

import (
	"bytes"
	"context"
	"errors"
	"ksef/core/invoicedb"
	"ksef/core/ksef"
	"ksef/internal/client/v2/invoices/downloader/metadata"
	"ksef/internal/logging"

	types "ksef/internal/client/v2/types/invoices"
)

var (
	errUnableToDownloadInvoice     = errors.New("unable to download invoice")
	errProcessingDownloadedInvoice = errors.New("error processing downloaded invoice")
)

func (sd *SimpleDownloader) DownloadInvoices(
	ctx context.Context,
	metadataPage *types.InvoiceMetadataResponse,
	subjectType ksef.SubjectType,
	invoiceReady func(subjectType ksef.SubjectType, invoiceMeta ksef.InvoiceMetadata, content bytes.Buffer) error,
) (err error) {
	var (
		finished             bool
		page                 int
		invoiceContentBuffer bytes.Buffer
	)

	for !finished {
		// fetch all of the invoices from the page.
		for _, invoice := range metadataPage.Invoices {
			yearMonth := invoicedb.YearMonth{
				Year:  sd.params.StartDate.Year(),
				Month: int(sd.params.StartDate.Month()),
			}
			exists, err := sd.hashChecker.ContainsHash(ctx, yearMonth, invoice.Checksum())
			if err != nil {
				return err
			}
			if exists {
				logging.DownloadLogger.Info("Ta faktura już została pobrana", "hash", invoice.Checksum())
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
