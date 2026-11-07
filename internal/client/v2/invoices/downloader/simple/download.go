package simple

import (
	"bytes"
	"context"
	"errors"
	"ksef/core/downloader"
	"ksef/core/ksef"
	"ksef/internal/client/v2/invoices/downloader/metadata"
	"ksef/internal/logging"
)

var errFetchingInvoices = errors.New("error fetching invoices")

func (sd *SimpleDownloader) Download(
	ctx context.Context,
	params downloader.Config,
	invoiceReady func(
		subjectType ksef.SubjectType,
		invoice ksef.InvoiceMetadata,
		content bytes.Buffer,
	) error,
) error {
	sd.params = params
	for _, subjectType := range params.SubjectTypes {
		logger := logging.DownloadLogger.With("subjectType", subjectType)
		logger.Debug("fetch invoices list", "start timestamp", params.StartDate)

		metadataPage, err := metadata.InvoicesMetadataPage(
			ctx,
			sd.httpClient,
			subjectType,
			params,
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
