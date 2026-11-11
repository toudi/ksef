package client

import "ksef/internal/client/v2/auth/validator"

type ClientInitializer func(*Client)

func WithPDFRenderer(renderer InvoicePDFRenderer) ClientInitializer {
	return func(c *Client) {
		c.pdfRenderer = renderer
	}
}

// WithAuthValidator sets the challenge validator used to build the token
// manager. A nil validator yields a restore-only manager (no challenge flow).
func WithAuthValidator(v validator.AuthChallengeValidator) ClientInitializer {
	return func(c *Client) {
		c.authValidator = v
	}
}

// WithTokenProvider overrides the internally constructed token manager. It is
// primarily a testing seam; when set, no token manager is constructed.
func WithTokenProvider(provider TokenProvider) ClientInitializer {
	return func(c *Client) {
		c.tokenProvider = provider
	}
}
