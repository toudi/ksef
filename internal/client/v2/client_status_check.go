package v2

import (
	"context"
	"ksef/internal/client/v2/session/status"
	"ksef/internal/client/v2/upo"
	"ksef/internal/logging"
)

// TODO:
// c.GetToken() - blokująca funkcja z mutexem.
// dzięki temu będzie można rozpocząć pobieranie i blokować do momentu uzyskania pierwszego
// (i każdego kolejnego) tokenu sesyjnego.

func (c *APIClient) UploadSessionsStatusCheck(ctx context.Context, upoDownloaderParams upo.UPODownloaderParams) error {
	httpClient := c.authenticatedHTTPClient()

	upoDownloader := upo.NewDownloader(httpClient, upoDownloaderParams)

	// let's iterate through all of upload sessions and skip these that were already processed
	for uploadSessionId, uploadSessionStatus := range c.registry.UploadSessions {
		if uploadSessionStatus.Processed {
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
		}

		c.registry.MarkUploadSessionProcessed(uploadSessionId)
	}

	return c.registry.Save("")
}
