package client

import (
	"time"

	"ksef/internal/http"

	ratelimitsclient "ksef/internal/client/v2/rate-limits"
)

const defaultTokenTimeout = 30 * time.Second

// authenticatedHTTPClient builds (once) a client that retrieves a bearer token
// before every request and applies the configured rate-limit/retry policy.
//
// Starting the token manager happens here, lazily: a command that never issues
// an authenticated request never starts the auth worker.
func (c *Client) authenticatedHTTPClient() *http.Client {
	c.authOnce.Do(func() {
		if runner, ok := c.tokenProvider.(tokenRunner); ok {
			go runner.Run()
		}

		tokenTimeout := c.httpConfig.TokenTimeout
		if tokenTimeout <= 0 {
			tokenTimeout = defaultTokenTimeout
		}

		c.authedHTTPClient = http.NewClient(
			c.httpClient.Base,
			http.WithAuthTokenRetriever(func(_ ...time.Duration) (string, error) {
				return c.tokenProvider.GetAuthorizationToken(tokenTimeout)
			}),
			http.WithRateLimitsDiscovery(ratelimitsclient.DiscoverRateLimits),
			http.WithRateLimitsState(c.httpConfig.RateLimitsState),
			http.WithMaxRetries(c.httpConfig.MaxRetries),
		)
	})
	return c.authedHTTPClient
}
