package simple

import (
	"bytes"
	"context"
	"errors"
	downloadertypes "ksef/internal/client/v2/invoices/downloader/types"
	"ksef/internal/client/v2/types/invoices"
	"ksef/internal/logging"
)

var errFetchingInvoices = errors.New("error fetching invoices")

func (sd *simpleDownloader) Download(
	ctx context.Context,
	invoiceReady func(
		subjectType invoices.SubjectType,
		invoice invoices.InvoiceMetadata,
		content bytes.Buffer,
	) error,
) (err error) {
	var req downloadertypes.InvoiceListRequest

	for _, subject := range sd.params.SubjectTypes {
		logger := logging.DownloadLogger.With("subjectType", subject)
		logger.Debug("fetch invoices list", "start timestamp", sd.params.StartDate)
		req = downloadertypes.InvoiceListRequest{
			SubjectType: subject,
			DateRange: downloadertypes.DateRange{
				DateType: downloadertypes.DateRangeStorage,
				From:     sd.params.StartDate,
				To:       sd.params.EndDate,
			},
		}

		if err = sd.fetchInvoices(
			ctx,
			req,
			sd.params,
			invoiceReady,
		); err != nil {
			return errors.Join(errFetchingInvoices, err)
		}
	}

	return nil
}
