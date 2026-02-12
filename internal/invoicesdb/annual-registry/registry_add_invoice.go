package annualregistry

import (
	"errors"
	"ksef/internal/sei"
	"ksef/internal/utils"
)

var (
	errGeneratingContent                 = errors.New("error serializing invoice content")
	errExistingInvoiceHasBeenSentAlready = errors.New("cannot override invoice as it was already sent to KSeF")
)

func (r *Registry) AddInvoice(
	parsed *sei.ParsedInvoice,
	checksum string,
	storeContents bool,
) (err error) {
	var invoiceContents string
	if storeContents {
		invoiceContents, err = utils.Base64GZippedCBor(parsed.Invoice, 80)
		if err != nil {
			return errors.Join(errGeneratingContent, err)
		}
	}
	// check if we have to override an existing invoice or is it a mere append operation
	existingInvoice := r.GetByRefNo(parsed.Invoice.Number)
	if existingInvoice != nil {
		if existingInvoice.KSeFRefNo != "" {
			return errExistingInvoiceHasBeenSentAlready
		}
		existingInvoice.Checksum = checksum
		existingInvoice.Contents = invoiceContents
		existingInvoice.GenerationTime = parsed.Invoice.GenerationTime
	} else {
		r.invoices = append(
			r.invoices,
			&Invoice{
				RefNo:          parsed.Invoice.Number,
				Contents:       invoiceContents,
				Checksum:       checksum,
				GenerationTime: parsed.Invoice.GenerationTime,
			},
		)
	}
	return nil
}
