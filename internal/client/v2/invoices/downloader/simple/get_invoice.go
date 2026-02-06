package simple

import (
	"bytes"
	"context"
	"fmt"
	"ksef/internal/client/v2/types/invoices"
	"ksef/internal/http"
	baseHTTP "net/http"
)

const (
	endpointDownloadInvoice = "/api/v2/invoices/ksef/%s"
)

func (sd *simpleDownloader) downloadInvoice(
	ctx context.Context,
	invoiceMeta invoices.InvoiceMetadata,
	dest *bytes.Buffer,
) (err error) {
	_, err = sd.httpClient.Request(
		ctx, http.RequestConfig{
			DestWriter:     dest,
			ExpectedStatus: baseHTTP.StatusOK,
		}, fmt.Sprintf(endpointDownloadInvoice, invoiceMeta.KSeFNumber),
	)

	return err
}
