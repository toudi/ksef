package certsdb

import "context"

func (c *CertsDB) StoreMinistryCertificate(
	ctx context.Context,
	certData MinistryCertificateData,
) (bool, error) {
	return c.metadataStorage.StoreMinistryCertificate(ctx, certData)
}

func (c *CertsDB) AddUserCertificate(
	ctx context.Context,
	newCert UserCertificateData,
) (bool, error) {
	if err := newCert.Certificate.ValidateUserCertificate(); err != nil {
		return false, err
	}

	return c.metadataStorage.StoreUserCertificate(ctx, newCert)
}
