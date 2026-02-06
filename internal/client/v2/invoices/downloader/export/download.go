package export

import (
	"bytes"
	"context"
	"errors"
	"ksef/internal/certsdb"
	downloadertypes "ksef/internal/client/v2/invoices/downloader/types"
	"ksef/internal/client/v2/types/invoices"
	"ksef/internal/encryption"
	"ksef/internal/http"
	baseHTTP "net/http"
)

var errFetchingInvoicesWithExport = errors.New("error fetching invoices with export endpoint")

func (ed *exportDownloader) Download(
	ctx context.Context,
	invoiceReady func(
		subjectType invoices.SubjectType,
		invoice invoices.InvoiceMetadata,
		content bytes.Buffer,
	) error,
) (err error) {
	// first we have to prepare the export request
	for _, subjectType := range ed.params.SubjectTypes {
		cipher, err := encryption.CipherInit(32)
		if err != nil {
			return err
		}

		req := exportRequest{
			Filters: downloadertypes.InvoiceListRequest{
				SubjectType: subjectType,
				DateRange: downloadertypes.DateRange{
					DateType: downloadertypes.DateRangeStorage,
					From:     ed.params.StartDate,
					To:       ed.params.EndDate,
				},
			},
		}

		certificate, err := ed.certsDB.GetByUsage(certsdb.UsageSymmetricKeyEncryption, "")
		if err != nil {
			return err
		}
		req.Encryption, err = cipher.PrepareHTTPRequestPayload(certificate.Filename())
		if err != nil {
			return err
		}

		var resp exportResponse

		_, err = ed.httpClient.Request(
			ctx,
			http.RequestConfig{
				ContentType:     http.JSON,
				DestContentType: http.JSON,
				Dest:            &resp,
				Body:            req,
				Method:          baseHTTP.MethodPost,
				ExpectedStatus:  baseHTTP.StatusCreated,
			},
			endpointInvoicesExport,
		)
		if err != nil {
			return err
		}

		if err = ed.fetchInvoices(
			ctx,
			cipher,
			req,
			resp,
			invoiceReady,
		); err != nil {
			return errors.Join(errFetchingInvoicesWithExport, err)
		}

	}

	return nil
}
