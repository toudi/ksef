package certsdb

import "context"

func (c *CertsDB) AddNIPToCertificate(
	ctx context.Context,
	certId string,
	NIP string,
) error {
	return c.metadataStorage.AddNIPToCertificate(ctx, certId, NIP)
}

func (c *CertsDB) RemoveNIPFromCertificate(
	ctx context.Context,
	certId string,
	NIP string,
) error {
	return c.metadataStorage.RemoveNIPFromCertificate(ctx, certId, NIP)
}
