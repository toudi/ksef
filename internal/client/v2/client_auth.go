package v2

import (
	"context"
	"errors"
	"time"

	"ksef/internal/client/v2/auth"
)

var (
	errTimeoutWaitingForTokenEvent = errors.New("timeout waiting for token manager event loop")
	errTokenManagerNotInitialized  = errors.New("token manager is not initialized")
)

const tokenManagerShutdownTimeout = 30 * time.Second

func (c *APIClient) WaitForTokenManagerLoop() error {
	if c.tokenManager == nil {
		return errTokenManagerNotInitialized
	}

	c.StopTokenManager()

	return c.awaitTokenManagerLoop()
}

func (c *APIClient) AwaitTokenManagerLoop() error {
	if c.tokenManager == nil {
		return errTokenManagerNotInitialized
	}

	return c.awaitTokenManagerLoop()
}

func (c *APIClient) awaitTokenManagerLoop() error {
	ctx, cancel := context.WithTimeout(context.Background(), tokenManagerShutdownTimeout)
	defer cancel()

	select {
	case <-c.tokenManager.Done():
		return nil
	case <-ctx.Done():
		return errTimeoutWaitingForTokenEvent
	}
}

// ObtainToken blocks until session tokens are available (restored or redeemed).
func (c *APIClient) ObtainToken() error {
	ctx, cancel := context.WithTimeout(context.Background(), tokenManagerShutdownTimeout)
	defer cancel()

	return c.tokenManager.WaitReady(ctx)
}

// Logout stops the token manager's event loop.
func (c *APIClient) Logout() error {
	c.StopTokenManager()
	return nil
}

func (c *APIClient) BindNIPToPESEL(ctx context.Context, nip, pesel string) error {
	return auth.BindNIPToPESEL(ctx, c.httpClient, nip, pesel)
}

// StartTokenManager starts the auth event loop in the background. It is
// idempotent.
func (c *APIClient) StartTokenManager() error {
	if c.tokenManager == nil {
		return errTokenManagerNotInitialized
	}

	if c.tokenManagerStarted {
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	c.tokenManagerCancel = cancel
	go c.tokenManager.Run(ctx)
	c.tokenManagerStarted = true

	return nil
}

// StopTokenManager asks the event loop to finish. It is safe to call when the
// manager was never started, and it is idempotent.
func (c *APIClient) StopTokenManager() {
	if c.tokenManagerCancel != nil {
		c.tokenManagerCancel()
		c.tokenManagerCancel = nil
	}
	c.tokenManagerStarted = false
}
