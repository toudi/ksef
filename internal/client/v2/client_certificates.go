package v2

import (
	"ksef/internal/client/v2/certificates"
	"ksef/internal/runtime"
)

func (c *APIClient) Certificates(env runtime.Gateway) (*certificates.Manager, error) {
	if c.certificates == nil {
		c.certificates = certificates.NewManager(c.authenticatedHTTPClient(), c.certsDB, env)
	}

	return c.certificates, nil
}
