package v2

import (
	"context"
	"ksef/internal/config"
	httpClient "ksef/internal/http"
	"ksef/internal/logging"
	registryPkg "ksef/internal/registry"
	"ksef/internal/sei/api/client/v2/auth"
	"ksef/internal/sei/api/client/v2/security"
	"ksef/internal/sei/api/client/v2/session/interactive"
)

type APIClient struct {
	auth       *auth.AuthHandler
	Auth       *auth.Manager
	apiConfig  config.APIConfig
	httpClient httpClient.Client
	ctx        context.Context
	// for uploading sessions
	registry *registryPkg.InvoiceRegistry
}

func NewClient(ctx context.Context, cfg config.Config, env config.APIEnvironment, options ...func(c *APIClient)) (*APIClient, error) {
	logging.SeiLogger.Info("klient KSeF v2 - start programu")

	apiConfig := cfg.APIConfig(env)
	httpClient := httpClient.Client{Base: "https://" + apiConfig.Host}

	client := &APIClient{
		ctx:        ctx,
		httpClient: httpClient,
		apiConfig:  apiConfig,
	}

	for _, option := range options {
		option(client)
	}

	return client, nil
}

func (c *APIClient) DownloadCertificates(ctx context.Context) error {
	logging.SeiLogger.Debug("pobieranie certyfikatów klucza publicznego")

	return security.DownloadCertificates(
		ctx, c.httpClient, c.apiConfig,
	)
}

func (c *APIClient) InteractiveSession() (*interactive.Session, error) {
	collection, err := c.registry.InvoiceCollection()
	if err != nil {
		return nil, err
	}

	issuerNip := collection.Issuer

	// WARNING: for now, the challenge validator is forced to ksefToken
	challengeValidator := auth.NewKsefTokenAuthValidator(c.apiConfig, issuerNip)

	c.auth = auth.NewAuthHandler(c.httpClient, c.apiConfig, auth.WithChallengeValidator(
		challengeValidator,
	))
	c.Auth = auth.NewManager(c.auth)
	c.httpClient.AuthTokenRetrieverFunc = c.auth.GetAuthorizationToken

	return interactive.NewSession(c.httpClient, c.registry), nil
}

func WithRegistry(registry *registryPkg.InvoiceRegistry) func(client *APIClient) {
	return func(client *APIClient) {
		client.registry = registry
	}
}
