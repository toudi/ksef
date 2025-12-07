package annualregistry

import (
	"errors"
	"ksef/internal/sei"
	"ksef/internal/utils"
)

var (
	errGeneratingContent = errors.New("error serializing invoice content")
)

func (r *Registry) AddInvoice(
	parsed *sei.ParsedInvoice,
	fileName string,
	checksum string,
) error {
	invoiceContents, err := utils.Base64GZippedCBor(parsed.Invoice, 80)
	if err != nil {
		return errors.Join(errGeneratingContent, err)
	}
	r.invoices = append(
		r.invoices,
		&Invoice{
			RefNo:          parsed.Invoice.Number,
			Contents:       invoiceContents,
			Checksum:       checksum,
			GenerationTime: parsed.Invoice.GenerationTime,
		},
	)
	return nil
}
