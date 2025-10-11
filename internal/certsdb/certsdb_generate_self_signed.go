package certsdb

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"os"
	"path"
	"time"
)

// based on https://go.dev/src/crypto/tls/generate_cert.go

func (c *CertificatesDB) GenerateSelfSignedCert(subject pkix.Name, basename string) error {
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

	certOutDir := path.Dir(certificatesDBFile)
	certOutFilename := path.Join(certOutDir, fmt.Sprintf("%s-cert.pem", basename))
	certOut, err := os.Create(certOutFilename)
	if err != nil {
		log.Fatalf("Failed to open cert.pem for writing: %v", err)
	}
	if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); err != nil {
		log.Fatalf("Failed to write data to cert.pem: %v", err)
	}
	if err := certOut.Close(); err != nil {
		log.Fatalf("Error closing cert.pem: %v", err)
	}

	keyOutFilename := path.Join(certOutDir, fmt.Sprintf("%s-key.pem", basename))
	keyOut, err := os.OpenFile(keyOutFilename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Failed to open key.pem for writing: %v", err)
	}
	// for RSA:
	// privBytes, err := x509.MarshalPKCS8PrivateKey(private)
	privBytes, err := x509.MarshalECPrivateKey(private)
	if err != nil {
		log.Fatalf("Unable to marshal private key: %v", err)
	}
	if err := pem.Encode(keyOut, &pem.Block{Type: "EC PRIVATE KEY", Bytes: privBytes}); err != nil {
		log.Fatalf("Failed to write data to key.pem: %v", err)
	}
	if err := keyOut.Close(); err != nil {
		log.Fatalf("Error closing key.pem: %v", err)
	}
	return nil
}
