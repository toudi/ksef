package certsdb

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"ksef/core/ksef"
	"log"
	"math/big"
	"time"
)

func encodePrivateKeyBytes(privateKey *ecdsa.PrivateKey) ([]byte, error) {
	privateKeyBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return nil, err
	}
	var pemBytesBuffer bytes.Buffer
	if err = pem.Encode(&pemBytesBuffer, &pem.Block{Type: "PRIVATE KEY", Bytes: privateKeyBytes}); err != nil {
		return nil, err
	}

	return pemBytesBuffer.Bytes(), nil
}

// based on https://go.dev/src/crypto/tls/generate_cert.go

func (c *CertsDB) RegenerateSelfSignedCert(ctx context.Context, subject pkix.Name) error {
	var err error

	validFor := time.Duration(10 * 365 * 24 * time.Hour) // 10 years
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

	issuedAt := time.Now().Truncate(24 * time.Hour)

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

	privateKeyBytes, err := encodePrivateKeyBytes(private)
	if err != nil {
		return err
	}

	return c.storage.UpsertSelfSignedCertificate(
		ctx,
		UserCertificateData{
			Certificate: Certificate{
				SelfSigned: true,
				CertificateHash: CertificateHash{
					Gateway:   ksef.KSeFGatewayTest,
					Usage:     []CertificateUsage{CertificateUsageAuthentication},
					ValidFrom: template.NotBefore,
					ValidTo:   template.NotAfter,
				},
			},
			PrivateKeyPEMBytes: privateKeyBytes,
			DERBytes:           derBytes,
		},
	)
}
