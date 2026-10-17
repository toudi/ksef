package certificates

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"ksef/core/certsdb"
	"ksef/core/ksef"

	kr "ksef/core/keyring"
)

type EnrollmentType string

func (m *Manager) PrepareEnrollmentCSR(data *EnrollmentsData, usage certsdb.CertificateUsage, nip string, keyring kr.Keyring) error {
	template := &x509.CertificateRequest{
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

	ctx := context.Background()

	gateway, err := ksef.ParseGateway(m.env)
	if err != nil {
		return err
	}

	return m.certsDB.AddCSRRequest(
		ctx,
		certsdb.CSRRequest{
			Gateway:    gateway,
			PrivateKey: private,
			NIP:        nip,
			CSRData:    csrDerBytes,
			Usage:      usage,
		},
	)
}
