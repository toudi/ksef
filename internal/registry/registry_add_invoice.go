package registry

import (
	"errors"
	"ksef/internal/certsdb"
	"ksef/internal/client/v2/types/invoices"
	"ksef/internal/registry/types"
)

var (
	errCertificateMissingAndOfflineModeSelected = errors.New("nie wskazano certyfikatu dla faktury Offline")
)

func (r *InvoiceRegistry) AddInvoice(
	invoice invoices.InvoiceMetadata,
	checksum string,
	certificate *certsdb.Certificate,
) error {
	var regInvoice = types.Invoice{
		Metadata:            invoice.Metadata,
		ReferenceNumber:     invoice.InvoiceNumber,
		KSeFReferenceNumber: invoice.KSeFNumber,
		InvoiceType:         invoice.InvoiceType,
		IssueDate:           invoice.IssueDate,
		Checksum:            checksum,
		Offline:             invoice.Offline,
		SubjectFrom: types.InvoiceSubject{
			TIN: invoice.Seller.NIP,
		},
	}
	// let's generate qrcodes. the qrcode for the invoice itself is always the same:
	invoiceQRCode, err := GenerateInvoiceQRCode(r.Environment, regInvoice)
	if err != nil {
		return err
	}
	regInvoice.QRCodes.Invoice = invoiceQRCode

	// if it's an offline invoice, we have to generate qrcode for it
	if invoice.Offline {
		if certificate == nil {
			return errCertificateMissingAndOfflineModeSelected
		}
		certificateQRCode, err := GenerateCertificateQRCode(
			r.Environment,
			regInvoice,
			"Nip",
			*certificate,
		)
		if err != nil {
			return err
		}
		regInvoice.QRCodes.Certificate = certificateQRCode
	}
	r.Invoices = append(r.Invoices, regInvoice)
	r.seiRefNoIndex[invoice.KSeFNumber] = len(r.Invoices)
	r.checksumIndex[checksum] = len(r.Invoices)
	return nil
}
