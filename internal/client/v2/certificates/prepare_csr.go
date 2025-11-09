package certificates

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/base64"
	"ksef/internal/certsdb"
)

type EnrollmentType string

func (m *Manager) PrepareEnrollmentCSR(data *EnrollmentsData, usage certsdb.Usage, nip string) error {
	var template = &x509.CertificateRequest{
		Subject: pkix.Name{
			CommonName: data.CommonName,
			Country:    []string{data.CountryName},
		},
	}
	if data.SerialNumber != nil {
		template.Subject.ExtraNames = append(template.Subject.ExtraNames, pkix.AttributeTypeAndValue{
			Type:  asn1.ObjectIdentifier{2, 5, 4, 5},
			Value: *data.SerialNumber,
		})
	}
	if data.Surname != nil {
		template.Subject.ExtraNames = append(template.Subject.ExtraNames, pkix.AttributeTypeAndValue{
			Type:  asn1.ObjectIdentifier{2, 5, 4, 4},
			Value: *data.Surname,
		})
	}
	if data.GivenName != nil {
		template.Subject.ExtraNames = append(template.Subject.ExtraNames, pkix.AttributeTypeAndValue{
			Type:  asn1.ObjectIdentifier{2, 5, 4, 42},
			Value: *data.GivenName,
		})
	}
	if data.UniqueIdentifier != nil {
		template.Subject.ExtraNames = append(template.Subject.ExtraNames, pkix.AttributeTypeAndValue{
			Type:  asn1.ObjectIdentifier{2, 5, 4, 45},
			Value: *data.UniqueIdentifier,
		})
	}

	// let's prepare private key first
	private, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return err
	}
	csrDerBytes, err := x509.CreateCertificateRequest(rand.Reader, template, private)
	if err != nil {
		return err
	}

	return m.certsDB.AddCert(func(newCert *certsdb.Certificate) error {
		newCert.Environment = m.env
		var err error
		if err = newCert.SavePKey(private); err != nil {
			return err
		}
		newCert.NIP = nip
		newCert.CSRData = base64.StdEncoding.EncodeToString(csrDerBytes)
		newCert.Usage = append(newCert.Usage, usage)
		return nil
	})
}
