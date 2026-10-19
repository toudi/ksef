package certsdb

import "context"

func (c *CertsDB) UpdateCertificate(ctx context.Context, certId string, updateData CertificateUpdateData) error {
	if err := c.metadataStorage.UpdateCertificate(ctx, certId, updateData); err != nil {
		return err
	}

	// the metadata storage does not write raw bytes - the updated certificate
	// DER is persisted through the data storage
	if len(updateData.DERBytes) > 0 {
		return c.dataStorage.StoreCertificateBytes(ctx, certId, updateData.DERBytes)
	}

	return nil
}
