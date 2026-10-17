package certsdb

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"ksef/core/ksef"
)

type CSRRequest struct {
	PrivateKey *ecdsa.PrivateKey
	CSRData    []byte
	NIP        string
	Gateway    ksef.KSeFGateway
	Usage      CertificateUsage
}

func (c *CertsDB) AddCSRRequest(
	ctx context.Context,
	csrReq CSRRequest,
) error {
	// store the private key
	certificate := Certificate{
		CertificateHash: CertificateHash{
			Usage:   []CertificateUsage{csrReq.Usage},
			Gateway: csrReq.Gateway,
		},
		NIP:     []string{csrReq.NIP},
		CSRData: base64.StdEncoding.EncodeToString(csrReq.CSRData),
	}
	// a CSR request is a request for a user certificate, so it must satisfy
	// the user-certificate usage rules (government usages cannot be requested)
	if err := certificate.ValidateUserCertificate(); err != nil {
		return err
	}

	// store the private key
	privateKeyBytes, err := x509.MarshalPKCS8PrivateKey(csrReq.PrivateKey)
	if err != nil {
		return err
	}
	var pemBytesBuffer bytes.Buffer
	if err = pem.Encode(&pemBytesBuffer, &pem.Block{Type: "PRIVATE KEY", Bytes: privateKeyBytes}); err != nil {
		return err
	}
	privateKeyId, err := c.storage.StorePrivateKey(ctx, pemBytesBuffer.Bytes(), string(csrReq.Gateway))
	if err != nil {
		return err
	}

	certificate.PrivateKeyId = privateKeyId
	_, err = c.storage.StoreCSRRequest(ctx, certificate)
	return err
}
