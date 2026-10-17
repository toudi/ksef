package certificates

import (
	"context"
	"crypto/x509"
	"encoding/base64"
	"ksef/core/certsdb"
	"ksef/internal/logging"
)

func (m *Manager) SyncEnrollments(ctx context.Context) error {
	var err error
	// check if we have some awaiting csr's to send:
	awaitingSend, err := m.certsDB.FilterCertificates(ctx, func(c certsdb.Certificate) bool {
		return c.CSRData != "" && c.ReferenceNumber == ""
	})
	if err != nil {
		return err
	}

	for _, cert := range awaitingSend {
		logging.CertsDBLogger.With("typ", cert.UsageAsString()).Debug("wysyłam wniosek o certyfikat")
		resp, err := m.PushCSR(ctx, cert)
		if err != nil {
			return err
		}
		if err = m.certsDB.UpdateCertificate(
			ctx,
			cert.UID,
			certsdb.CertificateUpdateData{
				ReferenceNumber: &resp.ReferenceNumber,
			},
		); err != nil {
			return err
		}
	}

	// now check if we have some awaiting certs to download
	awaitingDownload, err := m.certsDB.FilterCertificates(ctx, func(c certsdb.Certificate) bool {
		return c.ReferenceNumber != "" && c.SerialNumber == ""
	})
	if err != nil {
		return err
	}

	for _, cert := range awaitingDownload {
		statusResp, err := m.GetEnrollmentStatus(ctx, cert)
		if err != nil {
			return err
		}
		if statusResp.Status.Code != enrollmentSuccess || statusResp.SerialNumber == nil {
			logging.CertsDBLogger.Debug("certyfikat nieprzetworzony, pomijam", "cert", cert.UID)
			continue
		}
		logging.CertsDBLogger.Debug("pobieram certyfikat")
		certResp, err := m.DownloadCertificate(ctx, *statusResp.SerialNumber)
		if err != nil {
			return err
		}
		certSerialNumber := *statusResp.SerialNumber
		derBytes, err := base64.StdEncoding.DecodeString(certResp.Certificates[0].Certificate)
		if err != nil {
			return err
		}
		certData, err := x509.ParseCertificate(derBytes)
		if err != nil {
			return err
		}
		if err = m.certsDB.UpdateCertificate(
			ctx,
			cert.UID,
			certsdb.CertificateUpdateData{
				SerialNumber: &certSerialNumber,
				ValidFrom:    &certData.NotBefore,
				ValidTo:      &certData.NotAfter,
				DERBytes:     derBytes,
			},
		); err != nil {
			return err
		}
	}

	return nil
}
