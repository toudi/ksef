package certsdb

import "context"

func (c *CertsDB) AddNIPToCertificate(
	ctx context.Context,
	certId string,
	NIP string,
) error {
	return c.storage.AddNIPToCertificate(ctx, certId, NIP)
}

func (c *CertsDB) RemoveNIPFromCertificate(
	ctx context.Context,
	certId string,
	NIP string,
) error {
	return c.storage.RemoveNIPFromCertificate(ctx, certId, NIP)
}
