package v2

import (
	"ksef/internal/sei/api/client/v2/session/batch"
	"ksef/internal/sei/api/client/v2/session/interactive"
)

func (c *APIClient) InteractiveSession() (*interactive.Session, error) {
	return interactive.NewSession(
		c.authenticatedHTTPClient(),
		c.registry,
	), nil
}

func (c *APIClient) BatchSession() (*batch.Session, error) {
	return batch.NewSession(
		c.authenticatedHTTPClient(),
		c.registry,
		c.apiConfig,
	), nil
}
