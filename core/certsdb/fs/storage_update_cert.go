package fs

import (
	"context"
	"encoding/pem"
	"ksef/core/certsdb"
	"os"
)

func (fs *fsCertsStorage) UpdateCertificate(ctx context.Context, certId string, updateData certsdb.CertificateUpdateData) error {
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
	if updateData.DERBytes != nil {
		if err = saveCertDERBytes(cert, updateData.DERBytes); err != nil {
			return err
		}
	}
	if updateData.ProfileName != nil {
		cert.ProfileName = *updateData.ProfileName
	}
	fs.replaceCertIndex(cert, currentHash)
	return nil
}

func saveCertDERBytes(cert *fsCertificate, derBytes []byte) error {
	certFile, err := os.Create(cert.Filename())
	if err != nil {
		return err
	}
	defer certFile.Close()

	return pem.Encode(certFile, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: derBytes,
	})
}
