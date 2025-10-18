package certsdb

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"log"
	"math/big"
	"time"
)

// based on https://go.dev/src/crypto/tls/generate_cert.go

func (c *CertificatesDB) GenerateSelfSignedCert(subject pkix.Name) error {
	var err error

	var validFor = time.Duration(10 * 365 * 24 * time.Hour) // 10 years
	// let's start by generating EC private key
	// _, priv, err = ed25519.GenerateKey(rand.Reader)
	// private, err := rsa.GenerateKey(rand.Reader, 2048)
	private, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	// public, private, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return err
	}

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return err
	}

	var issuedAt = time.Now().Truncate(24 * time.Hour)

	template := x509.Certificate{
		Subject:      subject,
		SerialNumber: serialNumber,
		NotBefore:    issuedAt,
		NotAfter:     issuedAt.Add(validFor),

		KeyUsage: x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		// ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA:                  true, // I presume this denotes the self-signed part
	}

	// for rsa:
	// derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &private.PublicKey, private)
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &private.PublicKey, private)
	if err != nil {
		log.Fatalf("Failed to create certificate: %v", err)
	}

	return c.Upsert(
		func(cert Certificate) bool {
			return cert.SelfSigned
		},
		func(newCert *Certificate) error {
			var err error
			newCert.SelfSigned = true
			newCert.Usage = []Usage{UsageAuthentication}
			newCert.ValidFrom = template.NotBefore
			newCert.ValidTo = template.NotAfter
			if err = newCert.SavePKey(private); err != nil {
				return err
			}
			if err = newCert.SaveCert(derBytes); err != nil {
				return err
			}

			return nil
		},
	)
}
