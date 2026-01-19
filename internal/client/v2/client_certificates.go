package v2

import (
	"ksef/internal/client/v2/certificates"
)

func (c *APIClient) Certificates(envId string) (*certificates.Manager, error) {
	if c.certificates == nil {
		c.certificates = certificates.NewManager(c.authenticatedHTTPClient(), c.certsDB, envId)
	}

	return c.certificates, nil
}
