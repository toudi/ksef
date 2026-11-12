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

	// authCtx is the root context the auth worker runs under. It is derived from
	// the context passed to Init (not from the context of whichever request
	// happens to start the worker first), so a request-scoped context can never
	// terminate the worker.
	authCtx     context.Context
	authCancel  context.CancelFunc
	authStarted bool
}

// authShutdownTimeout bounds how long Close waits for the auth worker.
const authShutdownTimeout = 5 * time.Second

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

	// Construct a separate context for the auth worker so that the requests don't accidentally
	// kill the worker itself (i.e. they will use separate contexts with separate timeouts)
	client.authCtx, client.authCancel = context.WithCancel(context.WithoutCancel(ctx))

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

// Close stops the background auth worker and waits (bounded) for it to finish.
// It is safe to call when no token manager was constructed (e.g. a token
// provider was injected) or when the worker was never started.
func (c *Client) Close() error {
	if c.authCancel == nil {
		return nil
	}
	c.authCancel()

	// Nothing to join if the worker was never started.
	if !c.authStarted {
		return nil
	}

	runner, ok := c.tokenProvider.(tokenRunner)
	if !ok {
		return nil
	}

	select {
	case <-runner.Done():
		return nil
	case <-time.After(authShutdownTimeout):
		return nil
	}
}
