package certificates

import (
	"context"
	"crypto/x509"
	"encoding/base64"
	"ksef/internal/certsdb"
	"ksef/internal/logging"
)

func (m *Manager) SyncEnrollments(ctx context.Context) error {
	var err error
	// check if we have some awaiting csr's to send:
	if err = m.certsDB.ForEach(
		func(cert certsdb.Certificate) bool {
			return cert.CSRData != "" && cert.ReferenceNumber == ""
		},
		func(cert *certsdb.Certificate) error {
			logging.CertsDBLogger.With("typ", cert.UsageAsString()).Debug("wysyłam wniosek o certyfikat")
			resp, err := m.PushCSR(ctx, cert)
			if err != nil {
				return err
			}
			cert.ReferenceNumber = resp.ReferenceNumber
			return nil
		},
	); err != nil {
		return err
	}
	// now check if we have some awaiting certs to download
	return m.certsDB.ForEach(
		func(cert certsdb.Certificate) bool {
			return cert.ReferenceNumber != "" && cert.SerialNumber == ""
		},
		func(cert *certsdb.Certificate) error {
			statusResp, err := m.GetEnrollmentStatus(ctx, cert)
			if err != nil {
				return err
			}
			if statusResp.Status.Code == enrollmentSuccess {
				// we're ready to download
				logging.CertsDBLogger.Debug("pobieram certyfikat")
				certResp, err := m.DownloadCertificate(ctx, *statusResp.SerialNumber)
				if err != nil {
					return err
				}
				cert.SerialNumber = *statusResp.SerialNumber
				derBytes, err := base64.StdEncoding.DecodeString(certResp.Certificates[0].Certificate)
				if err != nil {
					return err
				}
				certData, err := x509.ParseCertificate(derBytes)
				if err != nil {
					return err
				}
				cert.ValidFrom = certData.NotBefore
				cert.ValidTo = certData.NotAfter
				return cert.SaveCert(derBytes)
			}
			return nil
		},
	)
}
