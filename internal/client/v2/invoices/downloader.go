package invoices

import (
	"bytes"
	"context"
	"errors"
	"ksef/internal/client/v2/types/invoices"
	"ksef/internal/http"
	monthlyregistry "ksef/internal/invoicesdb/monthly-registry"
	"ksef/internal/logging"
)

var errFetchingInvoices = errors.New("error fetching invoices")

type InvoiceDownloader struct {
	httpClient *http.Client
	registry   *monthlyregistry.Registry
	params     invoices.DownloadParams
}

func NewInvoiceDownloader(
	httpClient *http.Client,
	downloadParams invoices.DownloadParams,
	registry *monthlyregistry.Registry,
) *InvoiceDownloader {
	return &InvoiceDownloader{
		httpClient: httpClient,
		registry:   registry,
		params:     downloadParams,
	}
}

func (d *InvoiceDownloader) Download(
	ctx context.Context,
	invoiceReady func(subjectType invoices.SubjectType, invoice invoices.InvoiceMetadata, content bytes.Buffer) error,
) (err error) {
	// startTimestamp := d.registry.SyncParams.LastTimestamp
	var req InvoiceMetadataRequest

	for _, subject := range d.params.SubjectTypes {
		logger := logging.DownloadLogger.With("subjectType", subject)
		logger.Debug("fetch invoices list", "start timestamp", d.params.StartDate)
		req = InvoiceMetadataRequest{
			SubjectType: subject,
			DateRange: DateRange{
				DateType: DateRangeStorage,
				From:     d.params.StartDate,
				To:       d.params.EndDate,
			},
		}

		if err = d.fetchInvoices(
			ctx,
			req,
			d.params,
			invoiceReady,
		); err != nil {
			return errors.Join(errFetchingInvoices, err)
		}
	}

	return nil
}
