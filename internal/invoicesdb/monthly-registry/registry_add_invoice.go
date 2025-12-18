package monthlyregistry

import (
	"encoding/base64"
	"errors"
	"ksef/internal/certsdb"
	"ksef/internal/client/v2/types/invoices"
	"ksef/internal/runtime"
	"ksef/internal/sei"
	"ksef/internal/utils"
)

var (
	errTryingToOverwritePushedInvoice = errors.New("you're trying to overwrite an invoice which was already pushed to KSeF")
	errUnableToGetCertificate         = errors.New("unable to get offline certificate")
	errUnableToGenerateQRCode         = errors.New("unable to generate QR Code")
)

func (r *Registry) AddInvoice(
	inv *sei.ParsedInvoice,
	invoiceType InvoiceType,
	checksum string,
) error {
	// check if the invoice with this checksum already exists
	if r.ContainsHash(checksum) {
		return nil
	}

	gateway := runtime.GetGateway(r.vip)

	// ok, so the invoice with the given checksum does not exist. let's check if
	// the invoice with the given number exists. if so - is the user trying to
	// overwrite it ?
	if invoice := r.getInvoiceByRefNo(inv.Invoice.Number); invoice != nil {
		// this is indeed the case. we can only allow it if the original invoice
		// was not pushed to KSeF yet:
		if invoice.KSeFRefNo != "" {
			return errTryingToOverwritePushedInvoice
		}
	}

	// we seem to be in the clear.
	invoice := &Invoice{
		RefNo:    inv.Invoice.Number,
		Checksum: checksum,
		Offline:  inv.Invoice.KSeFFlags.Offline,
		Type:     invoiceType,
	}
	var err error

	if invoice.QRCodes.Invoice, err = invoice.generateInvoiceQRCode(gateway, inv); err != nil {
		return errors.Join(errUnableToGenerateQRCode, err)
	}

	if inv.Invoice.KSeFFlags.Offline {
		// let's prepare the offline qrcode
		var certificate certsdb.Certificate
		if certificate, err = r.certsDB.GetByUsage(
			certsdb.UsageOffline,
			inv.Invoice.Issuer.NIP,
		); err != nil {
			return errors.Join(errUnableToGetCertificate, err)
		}

		if invoice.QRCodes.Offline, err = invoice.generateCertificateQRCode(
			gateway,
			inv,
			certificate,
		); err != nil {
			return errors.Join(errUnableToGenerateQRCode, err)
		}
	}

	r.Invoices = append(
		r.Invoices,
		invoice,
	)

	return nil
}

func (r *Registry) AddReceivedInvoice(ksefInvoice invoices.InvoiceMetadata, subjectType invoices.SubjectType, gateway runtime.Gateway) (err error) {
	var checksumBytes []byte
	checksumBytes, err = base64.StdEncoding.DecodeString(ksefInvoice.InvoiceHashBase64)
	if err != nil {
		return err
	}

	var invoice = &Invoice{
		RefNo:     ksefInvoice.InvoiceNumber,
		KSeFRefNo: ksefInvoice.KSeFNumber,
		Checksum:  utils.Base64ToHex(ksefInvoice.Checksum()),
		QRCodes: InvoiceQRCodes{
			Invoice: generateInvoiceQRCodeInner(
				string(gateway),
				ksefInvoice.Seller.NIP,
				ksefInvoice.IssueTime(),
				checksumBytes,
			),
		},
		Type: ksefSubjectTypeToRegistryInvoiceType[subjectType],
	}

	r.Invoices = append(
		r.Invoices,
		invoice,
	)

	return nil
}
