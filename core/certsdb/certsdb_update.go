package certsdb

import "context"

func (c *CertsDB) UpdateCertificate(ctx context.Context, certId string, updateData CertificateUpdateData) error {
	return c.storage.UpdateCertificate(ctx, certId, updateData)
}
