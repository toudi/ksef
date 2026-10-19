package fs

import (
	"context"
	"ksef/core/certsdb"
)

func (fs *FSCertsStorage) UpdateCertificate(ctx context.Context, certId string, updateData certsdb.CertificateUpdateData) error {
	cert, err := fs.getCertByID(certId)
	if err != nil {
		return err
	}
	currentHash := cert.Hash()
	if updateData.SerialNumber != nil {
		cert.SerialNumber = *updateData.SerialNumber
	}
	if updateData.ReferenceNumber != nil {
		cert.ReferenceNumber = *updateData.ReferenceNumber
	}
	if updateData.ValidFrom != nil {
		cert.ValidFrom = *updateData.ValidFrom
	}
	if updateData.ValidTo != nil {
		cert.ValidTo = *updateData.ValidTo
	}
	// the DERBytes are written by the data storage (the facade does it)
	if updateData.ProfileName != nil {
		cert.ProfileName = *updateData.ProfileName
	}
	fs.replaceCertIndex(cert, currentHash)
	return nil
}
