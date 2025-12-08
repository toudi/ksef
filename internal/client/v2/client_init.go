package v2

import (
	"context"
	"ksef/internal/certsdb"
	"ksef/internal/client/v2/auth"
	"ksef/internal/client/v2/auth/validator"
	"ksef/internal/client/v2/certificates"
	"ksef/internal/http"
	httpClient "ksef/internal/http"
	"ksef/internal/logging"
	registryPkg "ksef/internal/registry"
	"ksef/internal/runtime"

	"github.com/spf13/viper"
)

type APIClient struct {
	tokenManager           *auth.TokenManager
	authChallengeValidator validator.AuthChallengeValidator
	httpClient             *httpClient.Client
	ctx                    context.Context
	// for uploading sessions
	registry     *registryPkg.InvoiceRegistry
	certificates *certificates.Manager
	certsDB      *certsdb.CertificatesDB
	// init options
	runTokenManager bool
}

type InitializerFunc func(c *APIClient)

func NewClient(ctx context.Context, vip *viper.Viper, options ...InitializerFunc) (*APIClient, error) {
	logging.SeiLogger.Info("klient KSeF v2 - start programu")

	httpClient := http.NewClient(string(runtime.GetGateway(vip)))

	client := &APIClient{
		ctx:             ctx,
		httpClient:      httpClient,
		runTokenManager: true,
	}

	for _, option := range options {
		option(client)
	}

	if client.authChallengeValidator != nil {
		var err error
		if client.tokenManager, err = auth.NewTokenManager(ctx, vip, client.authChallengeValidator); err != nil {
			return nil, err
		}
		if client.runTokenManager {
			go client.tokenManager.Run()
		}
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

func WithCertificatesDB(certsDB *certsdb.CertificatesDB) func(client *APIClient) {
	return func(client *APIClient) {
		client.certsDB = certsDB
	}
}

func WithoutTokenManager() func(client *APIClient) {
	return func(client *APIClient) {
		client.runTokenManager = false
	}
}
