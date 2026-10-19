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
	certificate := Certificate{
		CertificateHash: CertificateHash{
			Usage:   []CertificateUsage{csrReq.Usage},
			Gateway: csrReq.Gateway,
		},
		NIP:     []string{csrReq.NIP},
		CSRData: base64.StdEncoding.EncodeToString(csrReq.CSRData),
	}
	// a CSR request is a request for a user certificate, so it must satisfy
	// the user-certificate usage rules (you cannot request a CSR with
	// usage that's applicable only to ministry of finance (TokenEncryption / SymmetricKeyEncryption)
	if err := certificate.ValidateUserCertificate(); err != nil {
		return err
	}

	privateKeyBytes, err := x509.MarshalPKCS8PrivateKey(csrReq.PrivateKey)
	if err != nil {
		return err
	}
	var pemBytesBuffer bytes.Buffer
	if err = pem.Encode(&pemBytesBuffer, &pem.Block{Type: "PRIVATE KEY", Bytes: privateKeyBytes}); err != nil {
		return err
	}

	// the private key UID is assigned by the metadata storage
	_, err = c.metadataStorage.StoreCSRRequest(
		ctx,
		CSRRequestData{
			Certificate:        certificate,
			PrivateKeyPEMBytes: pemBytesBuffer.Bytes(),
		},
	)
	return err
}
