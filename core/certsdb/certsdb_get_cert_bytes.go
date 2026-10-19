package certsdb

import "context"

func (c *CertsDB) GetCertificateBytes(
	ctx context.Context,
	certId string,
) ([]byte, error) {
	return c.dataStorage.GetCertificateBytes(ctx, certId)
}
