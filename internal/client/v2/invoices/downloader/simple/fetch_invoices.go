package simple

import (
	"bytes"
	"context"
	"errors"
	downloadertypes "ksef/internal/client/v2/invoices/downloader/types"
	"ksef/internal/client/v2/types/invoices"
	types "ksef/internal/client/v2/types/invoices"
	"ksef/internal/http"
	baseHttp "net/http"
	"strconv"
)

var (
	errUnableToDownloadInvoice     = errors.New("unable to download invoice")
	errProcessingDownloadedInvoice = errors.New("error processing downloaded invoice")
)

func (sd *simpleDownloader) fetchInvoices(
	ctx context.Context,
	req downloadertypes.InvoiceListRequest,
	params invoices.DownloadParams,
	invoiceReady func(subjectType invoices.SubjectType, invoiceMeta invoices.InvoiceMetadata, content bytes.Buffer) error,
) (err error) {
	var (
		finished             bool
		page                 int
		resp                 types.InvoiceMetadataResponse
		invoiceContentBuffer bytes.Buffer
	)

	for !finished {
		_, err = sd.httpClient.Request(
			ctx,
			http.RequestConfig{
				Method: baseHttp.MethodPost,
				QueryParams: map[string]string{
					"pageOffset": strconv.Itoa(page),
					"pageSize":   strconv.Itoa(params.PageSize),
				},
				Body:            req,
				ContentType:     http.JSON,
				Dest:            &resp,
				DestContentType: http.JSON,
				ExpectedStatus:  baseHttp.StatusOK,
			},
			endpointInvoicesMetadata,
		)
		if err != nil {
			return err
		}

		for _, invoice := range resp.Invoices {
			if sd.registry.ContainsHash(invoice.Checksum()) {
				continue
			}

			// download invoice to a buffer
			invoiceContentBuffer.Reset()

			if err = sd.downloadInvoice(ctx, invoice, &invoiceContentBuffer); err != nil {
				return errors.Join(errUnableToDownloadInvoice, err)
			}

			if err = invoiceReady(req.SubjectType, invoice, invoiceContentBuffer); err != nil {
				return errors.Join(errProcessingDownloadedInvoice, err)
			}
		}

		finished = !resp.HasMore

		page += 1
	}

	return nil
}
