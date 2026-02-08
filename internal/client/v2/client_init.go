package v2

import (
	"context"
	"ksef/internal/certsdb"
	"ksef/internal/client/v2/auth"
	"ksef/internal/client/v2/auth/validator"
	"ksef/internal/client/v2/certificates"
	ratelimits "ksef/internal/client/v2/rate-limits"
	"ksef/internal/http"
	httpClient "ksef/internal/http"
	"ksef/internal/logging"
	"ksef/internal/runtime"

	"github.com/spf13/viper"
)

type APIClient struct {
	tokenManager           *auth.TokenManager
	authChallengeValidator validator.AuthChallengeValidator
	httpClient             *httpClient.Client
	authedHTTPClient       *httpClient.Client
	ctx                    context.Context
	// for uploading sessions
	certificates *certificates.Manager
	certsDB      *certsdb.CertificatesDB
	// init options
	runTokenManager     bool
	tokenManagerStarted bool
	vip                 *viper.Viper
}

type InitializerFunc func(c *APIClient)

func NewClient(ctx context.Context, vip *viper.Viper, options ...InitializerFunc) (*APIClient, error) {
	logging.SeiLogger.Info("klient KSeF v2 - start programu")
	environment := runtime.GetEnvironment(vip)

	httpClient := http.NewClient(environment.API)

	client := &APIClient{
		ctx:             ctx,
		httpClient:      httpClient,
		runTokenManager: true,
		vip:             vip,
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
			if err = client.StartTokenManager(); err != nil {
				return nil, err
			}
		}
	}

	return client, nil
}

func (c *APIClient) authenticatedHTTPClient() *httpClient.Client {
	// always return the same instance so that we can call setRateLimits on it and all the other
	// endpoints make use of that.
	if c.authedHTTPClient == nil {
		// create a copy of httpClient that will use token manager instance to retrieve the current session token
		c.authedHTTPClient = &httpClient.Client{
			Base:                   c.httpClient.Base,
			AuthTokenRetrieverFunc: c.tokenManager.GetAuthorizationToken,
			RateLimitsDiscoverFunc: ratelimits.DiscoverRateLimits,
			Vip:                    c.vip,
		}
	}
	return c.authedHTTPClient
}

func (c *APIClient) Close() {
	c.tokenManager.Stop()
	if c.authedHTTPClient != nil {
		if err := c.authedHTTPClient.SaveRateLimitsState(c.vip); err != nil {
			logging.SeiLogger.Error("unable to persist rate limits state", "err", err)
		}
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
