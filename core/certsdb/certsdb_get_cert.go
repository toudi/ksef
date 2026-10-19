package certsdb

import "context"

// Certificate performs a single certificate lookup by UID.
func (c *CertsDB) GetCertificate(ctx context.Context, certId string) (*Certificate, error) {
	return c.metadataStorage.GetCertificate(ctx, certId)
}
