package v2

import (
	"ksef/internal/environment"
	"ksef/internal/sei/api/client/v2/certificates"
)

func (c *APIClient) Certificates(env environment.Environment) (*certificates.Manager, error) {
	if c.certificates == nil {
		c.certificates = certificates.NewManager(c.authenticatedHTTPClient(), c.certsDB, env)
	}

	return c.certificates, nil
}
