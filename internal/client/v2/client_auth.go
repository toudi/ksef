package v2

import (
	"context"
	"errors"
	"ksef/internal/client/v2/auth"
	"ksef/internal/http"
	"time"
)

var (
	errTimeoutWaitingForTokenEvent = errors.New("timeout waiting for token manager event loop")
	errTokenManagerNotInitialized  = errors.New("token manager is not initialized")
)

func (c *APIClient) WaitForTokenManagerLoop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
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

func (c *APIClient) SetSessionTokens(tokens *auth.SessionTokens) {
	c.tokenManager.SetSessionTokens(tokens)
}

func (c *APIClient) Logout() error {
	c.tokenManager.Stop()
	return nil
}

func (c *APIClient) BindNIPToPESEL(ctx context.Context, nip, pesel string) error {
	return auth.BindNIPToPESEL(ctx, c.httpClient, nip, pesel)
}

func (c *APIClient) StartTokenManager() error {
	if c.tokenManager == nil {
		return errTokenManagerNotInitialized
	}

	go c.tokenManager.Run()
	return nil
}

// yeah, I should probably rewrite it at some point.
// don't have the time now
func (c *APIClient) GetAuthedHTTPClient() *http.Client {
	return c.authenticatedHTTPClient()
}
