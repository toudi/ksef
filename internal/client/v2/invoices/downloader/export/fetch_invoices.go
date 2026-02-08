package export

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	ratelimits "ksef/internal/client/v2/rate-limits"
	"ksef/internal/client/v2/types/invoices"
	"ksef/internal/encryption"
	"ksef/internal/http"
	"ksef/internal/logging"
	"ksef/internal/runtime"
	baseHTTP "net/http"
	"time"
)

var (
	errTimeoutWaitingForExport = errors.New("przekroczono czas oczekiwania na eksport faktur")
	errInvalidStatusRequest    = errors.New("błędne zapytanie o status")
)

func (ed *exportDownloader) fetchInvoices(
	ctx context.Context,
	cipher *encryption.Cipher,
	exportRequest exportRequest,
	exportResponse exportResponse,
	invoiceReady func(
		subjectType invoices.SubjectType,
		invoice invoices.InvoiceMetadata,
		content bytes.Buffer,
	) error,
) error {
	// first - wait until status is "ready"
	ctxTimeout, ctxTimeoutCancel := context.WithTimeout(ctx, 10*time.Minute)
	defer ctxTimeoutCancel()

	var exportStatus *exportStatusResponse
	var exportStatusReady bool
	var err error

	exportStatus, exportStatusReady, err = ed.pollForExportStatus(ctxTimeout, exportResponse.ReferenceNumber)
	if err != nil {
		return err
	}

	statusPoller := time.NewTicker(runtime.HttpPollWaitTime(ed.vip))
	defer statusPoller.Stop()

	for !exportStatusReady {
		select {
		case <-ctxTimeout.Done():
			logging.DownloadLogger.Error("Przekroczono czas oczekiwania na eksport")
			return errTimeoutWaitingForExport
		case <-statusPoller.C:
			logging.DownloadLogger.Info("Czekam na wygenerowanie eksportu faktur", "refNo", exportResponse.ReferenceNumber)
			exportStatus, exportStatusReady, err = ed.pollForExportStatus(ctx, exportResponse.ReferenceNumber)
			if err != nil {
				return err
			}
		}
	}

	if exportStatus.Package.InvoiceCount == 0 {
		logging.DownloadLogger.Info("Brak faktur w paczce.")
		return nil
	}

	logging.DownloadLogger.Info("Eksport faktur gotowy. Przystępuję do pobierania faktur")

	return ed.downloadAndExtract(ctx, cipher, exportRequest, exportStatus, invoiceReady)
}

func (ed *exportDownloader) pollForExportStatus(
	ctx context.Context, referenceNumber string,
) (*exportStatusResponse, bool, error) {
	var esResp exportStatusResponse
	_, err := ed.httpClient.Request(
		ctx,
		http.RequestConfig{
			Method:          baseHTTP.MethodGet,
			DestContentType: http.JSON,
			Dest:            &esResp,
			ExpectedStatus:  baseHTTP.StatusOK,
			OperationId:     ratelimits.OperationExportStatus,
		},
		fmt.Sprintf(endpointInvoicesExportStatus, referenceNumber),
	)
	if err != nil {
		return nil, false, err
	}
	if invalid, description := esResp.Invalid(); invalid {
		return nil, false, errors.Join(errInvalidStatusRequest, errors.New(description))
	}
	return &esResp, esResp.Ready(), nil
}
