package registry

import (
	"errors"
	"ksef/internal/certsdb"
	"ksef/internal/sei/api/client/v2/types/invoices"
)

var (
	errCertificateMissingAndOfflineModeSelected = errors.New("nie wskazano certyfikatu dla faktury Offline")
)

func (r *InvoiceRegistry) AddInvoice(
	invoice invoices.InvoiceMetadata,
	checksum string,
	certificate *certsdb.Certificate,
) error {
	var regInvoice = Invoice{
		ReferenceNumber:     invoice.InvoiceNumber,
		KSeFReferenceNumber: invoice.KSeFNumber,
		InvoiceType:         invoice.InvoiceType,
		IssueDate:           invoice.IssueDate,
		Checksum:            checksum,
		Offline:             invoice.Offline,
		SubjectFrom: InvoiceSubject{
			TIN: invoice.Seller.NIP,
		},
	}
	if invoice.Offline {
		// let's generate qrcodes. also, let's make sure that the certificate is not nil
		if certificate == nil {
			return errCertificateMissingAndOfflineModeSelected
		}
		invoiceQRCode, err := GenerateInvoiceQRCode(r.Environment, regInvoice)
		if err != nil {
			return err
		}
		regInvoice.QRCodes.Invoice = invoiceQRCode
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
