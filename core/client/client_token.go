package client

import (
	"context"
	"time"

	"ksef/core/ksef"
)

// TokenProvider is the minimal capability the authenticated HTTP client needs:
// a bearer token for the current session.
type TokenProvider interface {
	GetAuthorizationToken(timeout ...time.Duration) (string, error)
}

// tokenRunner is the optional lifecycle a concrete token manager exposes. A
// provider injected for tests need not implement it.
type tokenRunner interface {
	Run(ctx context.Context)
	Done() <-chan struct{}
}

// authSessionManager is the optional capability used by the auth-session
// commands.
type authSessionManager interface {
	GetAuthSessions(ctx context.Context) (*ksef.AuthSessions, error)
	Logout(sessionRefNo string) error
}
