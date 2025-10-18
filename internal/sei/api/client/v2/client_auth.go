package v2

import (
	"context"
	"errors"
	"ksef/internal/environment"
	"ksef/internal/sei/api/client/v2/auth"
	"time"
)

var (
	errTimeoutWaitingForTokenEvent = errors.New("timeout waiting for token manager event loop")
)

func (c *APIClient) WaitForTokenManagerLoop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	select {
	case <-c.tokenManager.Done():
		return nil
	case <-ctx.Done():
		return errTimeoutWaitingForTokenEvent
	}
}

func (c *APIClient) ObtainToken() error {
	_, err := c.tokenManager.GetAuthorizationToken(30 * time.Second)
	if err != nil {
		return err
	}

	return nil
}

func (c *APIClient) PersistTokens(env environment.Environment, nip string) error {
	return c.tokenManager.PersistTokens(env, nip)
}
func (c *APIClient) SetSessionTokens(tokens *auth.SessionTokens) {
	c.tokenManager.SetSessionTokens(tokens)
}

func (c *APIClient) Logout() error {
	if err := c.tokenManager.Logout(); err != nil {
		return err
	}

	c.tokenManager.Stop()
	return nil
}

func (c *APIClient) BindNIPToPESEL(ctx context.Context, nip, pesel string) error {
	return auth.BindNIPToPESEL(ctx, c.httpClient, nip, pesel)
}
