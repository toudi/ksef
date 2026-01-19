package monthlyregistry

import (
	"encoding/base64"
	"encoding/hex"
	"errors"
	"ksef/internal/certsdb"
	"ksef/internal/sei"
	"net/url"
	"time"
)

const (
	contextIdentifierNIP = "Nip"
)

func generateInvoiceQRCodeInner(qrcodeURL string, issuerNIP string, issued time.Time, checksumBytes []byte) string {
	qrcode, _ := url.JoinPath(
		qrcodeURL,
		"client-app",
		"invoice",
		issuerNIP,
		issued.Format("02-01-2006"),
		base64.URLEncoding.EncodeToString(checksumBytes),
	)

	return qrcode
}

// https://github.com/CIRFMF/ksef-docs/blob/main/kody-qr.md#1-kodi--weryfikacja-i-pobieranie-faktury
func (i *Invoice) generateInvoiceQRCode(qrcodeURL string, parsed *sei.ParsedInvoice) (string, error) {
	checksumBytes, err := hex.DecodeString(i.Checksum)
	if err != nil {
		return "", errors.Join(errors.New("error converting checksum from hex to bytes"), err)
	}
	qrcode, _ := url.JoinPath(
		qrcodeURL,
		"client-app",
		"invoice",
		parsed.Invoice.Issuer.NIP,
		parsed.Invoice.Issued.Format("02-01-2006"),
		base64.URLEncoding.EncodeToString(checksumBytes),
	)
	return qrcode, nil
}

// https://github.com/CIRFMF/ksef-docs/blob/main/kody-qr.md#2-kodii--weryfikacja-certyfikatu
func (i *Invoice) generateCertificateQRCode(
	qrcodeURL string,
	parsed *sei.ParsedInvoice,
	certificate certsdb.Certificate,
) (string, error) {
	checksumBytes, err := hex.DecodeString(i.Checksum)
	if err != nil {
		return "", err
	}
	ctxIdentValue := parsed.Invoice.Issuer.NIP
	issuerNIP := parsed.Invoice.Issuer.NIP
	// note: here we do *not* use the leading https://
	signingContent, _ := url.JoinPath(
		qrcodeURL,
		"client-app",
		"certificate",
		contextIdentifierNIP,
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
