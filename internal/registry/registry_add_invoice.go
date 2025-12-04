package registry

import (
	"errors"
	"ksef/internal/certsdb"
	"ksef/internal/client/v2/types/invoices"
	"ksef/internal/registry/types"
)

var (
	errCertificateMissingAndOfflineModeSelected = errors.New("nie wskazano certyfikatu dla faktury Offline")
	errConflictingRefsForSameChecksum           = errors.New("different invoice already registered for the same checksum")
)

func (r *InvoiceRegistry) AddInvoice(
	invoice invoices.InvoiceMetadata,
	checksum string,
	certificate *certsdb.Certificate,
) error {
	var index int
	var exists bool

	// sanity checks:
	// 1. we can replace existing invoice with overriding it's reference number, but only if it was not registered
	// in KSeF
	if index, exists = r.checksumIndex[checksum]; exists {
		if r.Invoices[index].ReferenceNumber != invoice.InvoiceNumber && r.Invoices[index].KSeFReferenceNumber != "" {
			return errConflictingRefsForSameChecksum
		}
	}
	// 2. we can replace existing invoice with a different content only if the original invoice was not registered
	// in KSeF
	if index, exists = r.refNoIndex[invoice.InvoiceNumber]; exists {
		if r.Invoices[index].Checksum != checksum && r.Invoices[index].KSeFReferenceNumber != "" {
			return errConflictingRefsForSameChecksum
		}
	}

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
	// see if we're overwriting an existing invoice:
	if index, exists = r.refNoIndex[invoice.InvoiceNumber]; exists {
		currentInvoice := r.Invoices[index]
		delete(r.checksumIndex, currentInvoice.Checksum)
		delete(r.seiRefNoIndex, currentInvoice.KSeFReferenceNumber)
		r.Invoices[index] = regInvoice
	} else {
		// it's a new invoice.
		index = len(r.Invoices)
		r.Invoices = append(r.Invoices, regInvoice)
	}
	if invoice.KSeFNumber != "" {
		r.seiRefNoIndex[invoice.KSeFNumber] = index
	}
	r.checksumIndex[checksum] = index
	r.refNoIndex[invoice.InvoiceNumber] = index
	return nil
}
