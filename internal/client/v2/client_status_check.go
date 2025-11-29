package v2

import (
	"context"
	"errors"
	"ksef/internal/client/v2/session/status"
	"ksef/internal/client/v2/upo"
	"ksef/internal/logging"
	"time"
)

var (
	errTimeoutReached = errors.New("przekroczono maksymalny czas oczekiwania na UPO")
)

func (c *APIClient) UploadSessionsStatusCheck(ctx context.Context, upoDownloaderParams upo.UPODownloaderParams) error {
	httpClient := c.authenticatedHTTPClient()

	upoDownloader := upo.NewDownloader(httpClient, upoDownloaderParams)

	var allSessions = len(c.registry.UploadSessions)
	var processedSessions int

	var cancel context.CancelFunc

	if upoDownloaderParams.Wait > 0 {
		ctx, cancel = context.WithTimeout(ctx, upoDownloaderParams.Wait)
		defer cancel()
	}

	var pollDuration = time.Duration(5 * time.Second)

	for processedSessions < allSessions {
		// let's iterate through all of upload sessions and skip these that were already processed
		for uploadSessionId, uploadSessionStatus := range c.registry.UploadSessions {
			if uploadSessionStatus.Processed {
				processedSessions += 1
				continue
			}

			statusResponse, err := status.CheckSessionStatus(ctx, httpClient, uploadSessionId)
			if err != nil {
				// log this rather than returning an error eagerly I suppose since there may be
				// more than a single session to check
				logging.SeiLogger.Error("błąd sprawdzania statusu sesji", "sessionId", uploadSessionId)
				continue
			}

			if statusResponse.FailedInvoiceCount > 0 {
				// we have to mark failed invoices in registry
				failedInvoiceList, err := status.GetFailedInvoiceList(ctx, httpClient, uploadSessionId)
				if err != nil {
					logging.SeiLogger.Error("błąd pobierania listy błędnie przetworzonych faktur", "error", err)
				}
				c.registry.MarkFailedInvoices(uploadSessionId, failedInvoiceList)
			}

			if statusResponse.Status.Code == status.SessionStatusProcessed {
				logging.SeiLogger.Info("sesja przetworzona pomyślnie. Przystępuję do pobierania UPO")
				if err = upoDownloader.Download(ctx, uploadSessionId, statusResponse.Upo.Pages, func(upoRefNo string) {
					c.registry.AddUPOToSession(uploadSessionId, upoRefNo)
				}); err != nil {
					logging.SeiLogger.Error("błąd pobierania UPO", "error", err)
				}
				processedSessions += 1
			}

			c.registry.MarkUploadSessionProcessed(uploadSessionId)
		}

		// forloop is finished. let's check if we can return early.
		// we can do that if:
		// - user does not want to wait
		// - we managed to check all of the sessions on the first attempt and there's no point waiting
		if upoDownloaderParams.Wait == 0 || processedSessions == allSessions {
			return c.registry.Save("")
		}

		select {
		case <-ctx.Done():
			return errTimeoutReached
		default:
			// user wants to wait and there's no timeout yet - let's wait 5 seconds before continuing
			processedSessions = 0
			logging.UPOLogger.Debug("czekam 5s na ponowne sprawdzenie")
			time.Sleep(pollDuration)
		}
	}

	return c.registry.Save("")
}
