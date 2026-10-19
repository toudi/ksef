package fs

import (
	"context"
	"encoding/pem"
	"os"
)

func (d *FSDataStorage) certificateFileName(certId string) string {
	return d.fileName(certId, ".pem")
}

func (d *FSDataStorage) GetCertificateBytes(ctx context.Context, certId string) ([]byte, error) {
	return os.ReadFile(d.certificateFileName(certId))
}

// StoreCertificateBytes writes the certificate as a PEM file (PEM-encoded DER).
func (d *FSDataStorage) StoreCertificateBytes(ctx context.Context, certId string, derBytes []byte) error {
	certFile, err := os.Create(d.certificateFileName(certId))
	if err != nil {
		return err
	}
	defer certFile.Close()

	return pem.Encode(certFile, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: derBytes,
	})
}
