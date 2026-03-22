package smart

import (
	"bytes"
	"context"
	"ksef/internal/client/v2/invoices/downloader/metadata"
	ratelimits "ksef/internal/client/v2/rate-limits"
	"ksef/internal/client/v2/types/invoices"
	"ksef/internal/logging"
	"time"
)

func (sd *smartDownloader) Download(
	ctx context.Context,
	invoiceReady func(
		subjectType invoices.SubjectType,
		invoice invoices.InvoiceMetadata,
		content bytes.Buffer,
	) error,
) error {
	perMinuteLimit, err := sd.httpClient.GetLimit(
		ratelimits.OperationInvoiceDownload,
		time.Minute,
	)
	if err != nil {
		return err
	}

	// we give an arbitrary threshold. So basically if the number of awaiting invoices can be
	// downloaded within 2 minutes we will do so with the simple downloader as invoking export
	// mode has it's own lifecycle that for small number of inoices may be longer than doing
	// it directly.
	threshold := 2 * perMinuteLimit

	for _, subjectType := range sd.params.SubjectTypes {
		invoicesMetadata, err := metadata.InvoicesMetadataPage(
			ctx,
			sd.httpClient,
			subjectType,
			sd.params,
			0,
		)
		if err != nil {
			return err
		}

		if len(invoicesMetadata.Invoices) == 0 {
			logging.DownloadLogger.Info("No invoices for download", "subjectType", subjectType)
			continue
		}

		if len(invoicesMetadata.Invoices) <= threshold {
			if err := sd.downloadWithSimple(ctx, subjectType, invoicesMetadata, invoiceReady); err != nil {
				return err
			}
		} else {
			if err := sd.downloadWithExport(ctx, invoiceReady); err != nil {
				return err
			}
		}
	}

	return nil
}
