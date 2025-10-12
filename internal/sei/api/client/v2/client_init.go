package v2

import (
	"context"
	"ksef/internal/config"
	"ksef/internal/environment"
	httpClient "ksef/internal/http"
	"ksef/internal/logging"
	registryPkg "ksef/internal/registry"
	"ksef/internal/sei/api/client/v2/auth"
	"ksef/internal/sei/api/client/v2/auth/validator"
)

type APIClient struct {
	tokenManager           *auth.TokenManager
	authChallengeValidator validator.AuthChallengeValidator
	apiConfig              config.APIConfig
	httpClient             *httpClient.Client
	ctx                    context.Context
	// for uploading sessions
	registry *registryPkg.InvoiceRegistry
}

type InitializerFunc func(c *APIClient)

func NewClient(ctx context.Context, cfg config.Config, env environment.Environment, options ...InitializerFunc) (*APIClient, error) {
	logging.SeiLogger.Info("klient KSeF v2 - start programu")

	apiConfig := cfg.APIConfig(env)
	httpClient := &httpClient.Client{Base: "https://" + apiConfig.Environment.Host}

	client := &APIClient{
		ctx:        ctx,
		httpClient: httpClient,
		apiConfig:  apiConfig,
	}

	for _, option := range options {
		option(client)
	}

	if client.authChallengeValidator != nil {
		client.tokenManager = auth.NewTokenManager(httpClient, client.authChallengeValidator)
		go client.tokenManager.Run()
	}

	return client, nil
}

func (c *APIClient) authenticatedHTTPClient() *httpClient.Client {
	// create a copy of httpClient that will use token manager instance to retrieve the current session token
	return &httpClient.Client{
		Base:                   c.httpClient.Base,
		AuthTokenRetrieverFunc: c.tokenManager.GetAuthorizationToken,
	}
}

func (c *APIClient) Close() {
	c.tokenManager.Logout()
	c.tokenManager.Stop()
}

func WithRegistry(registry *registryPkg.InvoiceRegistry) func(client *APIClient) {
	return func(client *APIClient) {
		client.registry = registry
	}
}

func WithAuthValidator(validator validator.AuthChallengeValidator) func(client *APIClient) {
	return func(client *APIClient) {
		client.authChallengeValidator = validator
	}
}
