package simple

import (
	"bytes"
	"context"
	"errors"
	"ksef/internal/client/v2/invoices/downloader/metadata"
	"ksef/internal/client/v2/types/invoices"
	"ksef/internal/logging"
)

var errFetchingInvoices = errors.New("error fetching invoices")

func (sd *SimpleDownloader) Download(
	ctx context.Context,
	invoiceReady func(
		subjectType invoices.SubjectType,
		invoice invoices.InvoiceMetadata,
		content bytes.Buffer,
	) error,
) error {
	for _, subjectType := range sd.params.SubjectTypes {
		logger := logging.DownloadLogger.With("subjectType", subjectType)
		logger.Debug("fetch invoices list", "start timestamp", sd.params.StartDate)

		metadataPage, err := metadata.InvoicesMetadataPage(
			ctx,
			sd.httpClient,
			subjectType,
			sd.params,
			0,
		)
		if err != nil {
			return err
		}

		if err = sd.DownloadInvoices(
			ctx,
			metadataPage,
			subjectType,
			invoiceReady,
		); err != nil {
			return err
		}
	}

	return nil
}
