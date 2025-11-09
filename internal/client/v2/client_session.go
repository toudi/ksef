package v2

import (
	"errors"
	"ksef/internal/client/v2/session/batch"
	"ksef/internal/client/v2/session/interactive"
)

var (
	errCertsDBNotDefined = errors.New("nie zainicjowano bazy certyfikatów")
)

func (c *APIClient) InteractiveSession() (*interactive.Session, error) {
	if c.certsDB == nil {
		return nil, errCertsDBNotDefined
	}
	return interactive.NewSession(
		c.authenticatedHTTPClient(),
		c.registry,
		c.certsDB,
	), nil
}

func (c *APIClient) BatchSession() (*batch.Session, error) {
	if c.certsDB == nil {
		return nil, errCertsDBNotDefined
	}
	return batch.NewSession(
		c.authenticatedHTTPClient(),
		c.registry,
		c.certsDB,
	), nil
}
