package v2

import (
	"context"
	"ksef/internal/config"
	httpClient "ksef/internal/http"
	"ksef/internal/logging"
	"ksef/internal/sei/api/client/v2/auth"
	"ksef/internal/sei/api/client/v2/security"
)

type APIClient struct {
	auth       *auth.AuthHandler
	apiConfig  config.APIConfig
	httpClient httpClient.Client
	ctx        context.Context
}

func NewClient(ctx context.Context, cfg config.Config, env config.APIEnvironment) (*APIClient, error) {
	logging.SeiLogger.Info("klient KSeF v2 - start programu")

	apiConfig := cfg.APIConfig(env)
	httpClient := httpClient.Client{Base: "https://" + apiConfig.Host}

	return &APIClient{
		auth:       auth.NewAuthHandler(httpClient, apiConfig),
		ctx:        ctx,
		httpClient: httpClient,
		apiConfig:  apiConfig,
	}, nil
}

func (c *APIClient) DownloadCertificates(ctx context.Context) error {
	logging.SeiLogger.Debug("pobieranie certyfikatów klucza publicznego")

	return security.DownloadCertificates(
		ctx, c.httpClient, c.apiConfig,
	)
}
