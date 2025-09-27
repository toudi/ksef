package v2

import (
	"context"
	"ksef/internal/logging"
	"ksef/internal/sei/api/client/v2/security"
)

func (c *APIClient) DownloadPubKeyCertificates(ctx context.Context) error {
	logging.SeiLogger.Debug("pobieranie certyfikatów klucza publicznego")

	return security.DownloadCertificates(
		ctx, c.httpClient, c.apiConfig,
	)
}
