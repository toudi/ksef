package simple

import (
	"bytes"
	"context"
	"fmt"
	baseHTTP "net/http"

	ratelimits "ksef/internal/client/v2/rate-limits"
	"ksef/internal/client/v2/types/invoices"
	"ksef/internal/http"
)

const (
	endpointDownloadInvoice = "/v2/invoices/ksef/%s"
)

func (sd *SimpleDownloader) downloadInvoice(
	ctx context.Context,
	invoiceMeta invoices.InvoiceMetadata,
	dest *bytes.Buffer,
) (err error) {
	_, err = sd.httpClient.Request(
		ctx, http.RequestConfig{
			DestWriter:     dest,
			ExpectedStatus: baseHTTP.StatusOK,
			OperationId:    ratelimits.OperationInvoiceDownload,
		}, fmt.Sprintf(endpointDownloadInvoice, invoiceMeta.KSeFNumber),
	)

	return err
}
