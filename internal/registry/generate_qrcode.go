package registry

import (
	"encoding/base64"
	"encoding/hex"
	"ksef/internal/certsdb"
	"ksef/internal/config"
	"ksef/internal/registry/types"
	"net/url"
)

// https://github.com/CIRFMF/ksef-docs/blob/main/kody-qr.md#1-kodi--weryfikacja-i-pobieranie-faktury
func GenerateInvoiceQRCode(env config.Gateway, invoice types.Invoice) (string, error) {
	checksumBytes, err := hex.DecodeString(invoice.Checksum)
	if err != nil {
		return "", err
	}
	qrcode, _ := url.JoinPath(
		"https://"+string(env),
		"client-app",
		"invoice",
		invoice.SubjectFrom.TIN,
		invoice.IssueDate,
		base64.URLEncoding.EncodeToString(checksumBytes),
	)
	return qrcode, nil
}

// https://github.com/CIRFMF/ksef-docs/blob/main/kody-qr.md#2-kodii--weryfikacja-certyfikatu
func GenerateCertificateQRCode(
	env config.Gateway, invoice types.Invoice,
	ctxIdent string,
	certificate certsdb.Certificate,
) (string, error) {
	checksumBytes, err := hex.DecodeString(invoice.Checksum)
	if err != nil {
		return "", err
	}
	ctxIdentValue := invoice.SubjectFrom.TIN
	issuerNIP := invoice.SubjectFrom.TIN
	// note: here we do *not* use the leading https://
	signingContent, _ := url.JoinPath(
		string(env),
		"client-app",
		"certificate",
		ctxIdent,
		ctxIdentValue,
		issuerNIP,
		certificate.SerialNumber,
		base64.RawURLEncoding.EncodeToString(checksumBytes),
	)
	// now that we have signing content we can sign it with certificate
	signature, err := certificate.SignContent([]byte(signingContent))
	if err != nil {
		return "", err
	}
	return "https://" + signingContent + "/" + base64.RawURLEncoding.EncodeToString(signature), nil
}
