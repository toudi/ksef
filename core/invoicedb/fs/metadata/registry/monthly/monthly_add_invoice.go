package monthly

import (
	"encoding/base64"
	"fmt"
	"ksef/core/invoicedb"
	"ksef/core/invoicedb/invoice"
	"ksef/core/ksef"
	"ksef/internal/utils"
	"net/url"
)

func (r *Registry) AddReceivedInvoice(
	subjectType ksef.SubjectType,
	meta ksef.InvoiceMetadata,
	nip string,
	qrCodeBaseURL string,
) (*invoicedb.StoredInvoice, error) {
	invoiceType := subjectTypeToInvoiceType[subjectType]
	ordNum := r.ordNum[invoiceType] + 1

	// Received invoices are the ones that are issued "for" `nip`.
	// The code here generates the verification QRCode and, frankly, that is
	// only to be fully backwards compatible. I can only imagine that when I first
	// wrote the code I thought of rendering a PDF of a received invoice and for that
	// I needed the QRCode. Right now, as I'm refactoring this, I am strongly considering
	// to hide this behaviour behind a feature flag and only generate the QRCodes for
	// the issued invoices.
	checksumBytes, checksumErr := base64.StdEncoding.DecodeString(meta.InvoiceHashBase64)
	invoiceVerificationURL := ""
	if qrCodeBaseURL != "" {
		if checksumErr != nil {
			return nil, fmt.Errorf("unable to decode invoice hash: %w", checksumErr)
		}
		invoiceVerificationURL = generateInvoiceVerificationQRCode(qrCodeBaseURL, meta.Seller.NIP, meta.IssueDate, checksumBytes)
	}

	invoice := &Invoice{
		RefNo:     meta.InvoiceNumber,
		KSeFRefNo: meta.KSeFNumber,
		Checksum:  meta.Checksum(),
		Type:      invoiceType,
		OrdNum:    ordNum,
		QRCodes:   InvoiceQRCodes{Invoice: invoiceVerificationURL},
		Issuer: &InvoiceIssuer{
			NIP:  meta.Seller.NIP,
			Name: meta.Seller.Name,
		},
	}

	r.state.Invoices = append(r.state.Invoices, invoice)
	r.ordNum[invoiceType] = ordNum
	r.checksumIndex[invoice.Checksum] = len(r.state.Invoices) - 1
	r.dirty = true

	return &invoicedb.StoredInvoice{
		NIP:       nip,
		ID:        invoice.ID(),
		Invoice:   ToCoreInvoice(invoice),
		YearMonth: r.yearMonth,
	}, nil
}

func generateInvoiceVerificationQRCode(qrCodeBaseURL string, issuerNIP string, issueDate string, checksumBytes []byte) string {
	parsed, err := utils.ParseTimeFromString(issueDate)
	if err != nil {
		return ""
	}

	qrcode, err := url.JoinPath(
		qrCodeBaseURL,
		"invoice",
		issuerNIP,
		parsed.Format("02-01-2006"),
		base64.URLEncoding.EncodeToString(checksumBytes),
	)
	if err != nil {
		return ""
	}

	return qrcode
}

func ToCoreInvoice(input *Invoice) *invoice.Invoice {
	return &invoice.Invoice{
		RefNo:        input.RefNo,
		KSeFRefNo:    input.KSeFRefNo,
		Checksum:     input.Checksum,
		Offline:      input.Offline,
		QRCodes:      invoice.InvoiceQRCodes(input.QRCodes),
		Type:         invoice.InvoiceType(input.Type),
		UploadErrors: input.UploadErrors,
		PrintoutData: input.PrintoutData,
		Issuer:       ToCoreIssuer(input.Issuer),
		OrdNum:       input.OrdNum,
		Annotations:  invoice.Annotations(input.Annotations),
	}
}

func ToCoreIssuer(input *InvoiceIssuer) *invoice.InvoiceIssuer {
	if input == nil {
		return nil
	}
	return &invoice.InvoiceIssuer{NIP: input.NIP, Name: input.Name}
}
