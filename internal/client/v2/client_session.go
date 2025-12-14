package v2

import (
	"errors"
	"ksef/internal/client/v2/session/batch"
	"ksef/internal/client/v2/session/interactive"
	"ksef/internal/client/v2/session/types"
)

var (
	errCertsDBNotDefined = errors.New("nie zainicjowano bazy certyfikatów")
)

func (c *APIClient) InteractiveSession() (types.UploadSession, error) {
	if c.certsDB == nil {
		return nil, errCertsDBNotDefined
	}
	return interactive.NewSession(
		c.authenticatedHTTPClient(),
		c.certsDB,
	), nil
}

func (c *APIClient) BatchSession(workDir string) (types.UploadSession, error) {
	if c.certsDB == nil {
		return nil, errCertsDBNotDefined
	}
	return batch.NewSession(
		c.authenticatedHTTPClient(),
		c.certsDB,
		workDir,
	), nil
}
