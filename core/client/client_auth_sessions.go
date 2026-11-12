package client

import (
	"context"
	"errors"

	"ksef/core/ksef"
)

var errAuthSessionsUnsupported = errors.New("the configured token provider does not support auth sessions")

func (c *Client) AuthSessions(ctx context.Context) (*ksef.AuthSessions, error) {
	manager, ok := c.tokenProvider.(authSessionManager)
	if !ok {
		return nil, errAuthSessionsUnsupported
	}

	return manager.GetAuthSessions(ctx)
}

func (c *Client) LogoutSession(refNo string) error {
	manager, ok := c.tokenProvider.(authSessionManager)
	if !ok {
		return errAuthSessionsUnsupported
	}

	return manager.Logout(refNo)
}
