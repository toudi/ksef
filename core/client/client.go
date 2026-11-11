package client

import (
	"context"
	"sync"
	"time"

	"ksef/core/certsdb"
	"ksef/core/invoicedb"
	"ksef/core/keyring"
	"ksef/core/ksef"
	"ksef/core/ratelimits"
	"ksef/internal/client/v2/auth"
	"ksef/internal/client/v2/auth/validator"
	"ksef/internal/http"
)

// InvoicePDFRenderer renders a PDF for an already-stored invoice.
type InvoicePDFRenderer interface {
	RenderInvoice(
		ctx context.Context,
		stored *invoicedb.StoredInvoice,
	) error
}

// HTTPConfig carries the HTTP-client tuning for authenticated requests. The
// rate-limit state storage is optional: nil disables persistence. Rate-limit
// *discovery* is not configurable and is wired internally.
type HTTPConfig struct {
	ksef.HTTPConfig
	RateLimitsState ratelimits.StateStorage
	// TokenTimeout bounds how long a request waits for the session token. Zero
	// falls back to defaultTokenTimeout.
	TokenTimeout time.Duration
}

type Client struct {
	certsDB     *certsdb.CertsDB
	invoiceDB   *invoicedb.InvoiceDB
	keyring     keyring.Keyring
	httpClient  *http.Client
	gateway     ksef.KSeFGatewayConfig
	pdfRenderer InvoicePDFRenderer

	httpConfig    HTTPConfig
	authValidator validator.AuthChallengeValidator

	tokenProvider    TokenProvider
	authedHTTPClient *http.Client
	authOnce         sync.Once
}

func Init(
	ctx context.Context,
	certsDB *certsdb.CertsDB,
	invoiceDB *invoicedb.InvoiceDB,
	keyring keyring.Keyring,
	gateway ksef.KSeFGatewayConfig,
	nip string,
	httpConfig HTTPConfig,
	initializers ...ClientInitializer,
) (*Client, error) {
	client := &Client{
		certsDB:    certsDB,
		invoiceDB:  invoiceDB,
		keyring:    keyring,
		httpClient: http.NewClient(gateway.APIBaseURL),
		gateway:    gateway,
		httpConfig: httpConfig,
	}
	for _, initializer := range initializers {
		initializer(client)
	}

	// Construct the token manager eagerly so validator initialization errors
	// surface from Init; it is only *started* lazily on the first authenticated
	// request (see authenticatedHTTPClient).
	if client.tokenProvider == nil {
		manager, err := auth.NewTokenManager(ctx, auth.Config{
			NIP:           nip,
			EnvironmentID: string(gateway.Gateway),
			HTTPClient:    client.httpClient,
		}, client.authValidator, keyring, certsDB)
		if err != nil {
			return nil, err
		}
		client.tokenProvider = manager
	}

	return client, nil
}

// Close stops the background auth worker. It is safe to call when no token
// manager was constructed (e.g. a token provider was injected).
func (c *Client) Close() error {
	if runner, ok := c.tokenProvider.(tokenRunner); ok {
		runner.Stop()
	}
	return nil
}
