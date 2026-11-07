package http

import (
	"ksef/core/ratelimits"
	"ksef/internal/interfaces"
)

type Option func(*Client)

// WithMaxRetries overrides the default number of request attempts.
func WithMaxRetries(maxRetries int) Option {
	return func(c *Client) {
		c.MaxRetries = maxRetries
	}
}

// WithAuthTokenRetriever sets the function used to obtain the bearer token before
// each request. Requests without a retriever are sent unauthenticated.
func WithAuthTokenRetriever(retriever interfaces.TokenRetrieverFunc) Option {
	return func(c *Client) {
		c.AuthTokenRetrieverFunc = retriever
	}
}

// WithRateLimitsDiscovery sets the function used to discover the rate limits once
// an auth token is available.
func WithRateLimitsDiscovery(discover interfaces.RateLimitsDiscoverFunc) Option {
	return func(c *Client) {
		c.RateLimitsDiscoverFunc = discover
	}
}

// WithRateLimitsState enables persisting the rate limiter state through the given
// storage. The implementation is chosen by the caller, so the client itself does
// not need to know anything about the configuration.
func WithRateLimitsState(storage ratelimits.StateStorage) Option {
	return func(c *Client) {
		c.RateLimitsState = storage
	}
}
