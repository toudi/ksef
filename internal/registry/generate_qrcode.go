package registry

import (
	"encoding/base64"
	"encoding/hex"
	"errors"
	"ksef/internal/certsdb"
	"ksef/internal/registry/types"
	"ksef/internal/runtime"
	"net/url"
	"time"
)

// https://github.com/CIRFMF/ksef-docs/blob/main/kody-qr.md#1-kodi--weryfikacja-i-pobieranie-faktury
func GenerateInvoiceQRCode(env runtime.Gateway, invoice types.Invoice) (string, error) {
	checksumBytes, err := hex.DecodeString(invoice.Checksum)
	if err != nil {
		return "", errors.Join(errors.New("error converting checksum from hex to bytes"), err)
	}
	// the date has to be in the DD-MM-YYYY format
	issueDate, err := time.Parse(time.DateOnly, invoice.IssueDate)
	if err != nil {
		return "", errors.Join(errors.New("unable to parse issue date"), err)
	}
	qrcode, _ := url.JoinPath(
		"https://"+string(env),
		"client-app",
		"invoice",
		invoice.SubjectFrom.TIN,
		issueDate.Format("02-01-2006"),
		base64.URLEncoding.EncodeToString(checksumBytes),
	)
	return qrcode, nil
}

// https://github.com/CIRFMF/ksef-docs/blob/main/kody-qr.md#2-kodii--weryfikacja-certyfikatu
func GenerateCertificateQRCode(
	env runtime.Gateway, invoice types.Invoice,
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
