package statuschecker

import (
	"context"
	"errors"
	"ksef/internal/client/v2/session/status"
	"ksef/internal/client/v2/session/types"
	"ksef/internal/logging"
	"ksef/internal/pdf"
	pdfconfig "ksef/internal/pdf/config"
	"time"
)

var (
	errUnableToCheckStatus     = errors.New("unable to check session status")
	errTimeoutWaitingForStatus = errors.New("timed out waiting for upload session status")
	errPrintingInvoice         = errors.New("error printing invoice to PDF")
	errDownloadingUPO          = errors.New("error downloading UPO")
	errUnableToGetPrinter      = errors.New("unable to create PDF printer")
)

func (c *StatusChecker) CheckSessions(ctx context.Context) error {
	for sessionId, sessionRegistry := range c.sessionIdToSessionRegistry {
		logging.InvoicesDBLogger.Debug("checking session status", "session id", sessionId)

		sessionStatus, err := c.checkSessionStatus(ctx, sessionId)
		if err != nil {
			return errors.Join(errUnableToCheckStatus, err)
		}

		sessionRegistry.Update(
			sessionStatus,
			c.invoiceHashToMonthlyRegistry,
			c.invoiceHashToAnnualRegistry,
		)

		// if the session is processed we can (conditionally) download UPO
		if sessionStatus.IsProcessed() {
			if c.cfg.InvoicePDF {
				for invoiceHash, registry := range c.invoiceHashToMonthlyRegistry {
					invoice := registry.GetInvoiceByChecksum(
						invoiceHash,
					)
					if invoice == nil {
						return errors.New("unable to find invoice")
					}
					if invoice.Offline {
						// offline invoice was rendered to PDF during import - no need to re-render
						continue
					}

					sessionInvoice, err := sessionStatus.GetInvoiceByChecksum(invoiceHash)
					if err != nil {
						return err
					}

					printer, err := pdf.GetInvoicePrinter(c.vip, pdfconfig.UsageSelector{
						Usage:        "invoice:issued",
						Participants: sessionInvoice.Participants(),
					})
					if err != nil {
						return errors.Join(errUnableToGetPrinter, err)
					}

					invoiceFilename := registry.InvoiceFilename(invoice)
					if err = printer.PrintInvoice(
						invoiceFilename.XML,
						invoiceFilename.PDF,
						invoice.GetPrintingMeta(),
					); err != nil {
						return errors.Join(errPrintingInvoice, err)
					}
				}
			}

			if c.cfg.UPODownloaderConfig.Enabled {
				logging.InvoicesDBLogger.Info("sesja przetworzona pomyślnie. pobieram UPO", "upo", sessionStatus.Status.Upo)
				if err = c.downloadUPO(ctx, sessionStatus); err != nil {
					return errors.Join(errDownloadingUPO, err)
				}
			}
		}
	}

	return nil
}

func (c *StatusChecker) checkSessionStatus(
	ctx context.Context,
	sessionId string,
) (*types.UploadSessionResult, error) {
	var err error
	var uploadResult *types.UploadSessionResult

	// do the initial status check
	uploadResult, err = c.checkUploadSessionResult(ctx, sessionId)
	if err != nil {
		return nil, errors.Join(errUnableToCheckStatus, err)
	}

	if uploadResult.IsProcessed() {
		// perfect - we can return.
		return uploadResult, nil
	}

	// does the user want to wait ?
	if c.cfg.Wait {
		uploadResult, err = c.waitForResult(ctx, uploadResult)
	}

	return uploadResult, err
}

func (c *StatusChecker) checkUploadSessionResult(
	ctx context.Context,
	sessionId string,
) (*types.UploadSessionResult, error) {
	statusResponse, err := c.checker.CheckSessionStatus(ctx, sessionId)
	if err != nil {
		return nil, err
	}
	var invoiceList []status.InvoiceInfo
	if statusResponse.IsProcessed() {
		// let's download invoices for this session
		invoiceList, err = c.checker.GetInvoiceList(ctx, sessionId)
		if err != nil {
			return nil, err
		}
	}
	return &types.UploadSessionResult{
		SessionID: sessionId,
		Invoices:  invoiceList,
		Status:    statusResponse,
	}, nil
}

func (c *StatusChecker) waitForResult(
	ctx context.Context,
	uploadResult *types.UploadSessionResult,
) (*types.UploadSessionResult, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, c.cfg.WaitTimeout)
	defer cancel()

	pollingTicker := time.NewTicker(5 * time.Second)
	defer pollingTicker.Stop()

	var finished bool = false
	var err error

	for !finished {
		select {
		case <-timeoutCtx.Done():
			return nil, errTimeoutWaitingForStatus
		case <-pollingTicker.C:
			uploadResult, err = c.checkUploadSessionResult(ctx, uploadResult.SessionID)
			if err != nil {
				return nil, errors.Join(errUnableToCheckStatus, err)
			}
			finished = uploadResult.IsProcessed()
		}
	}

	return uploadResult, nil
}
