package invoicesdb

import (
	annualregistry "ksef/internal/invoicesdb/annual-registry"
	"ksef/internal/sei"
)

func (idb *InvoicesDB) handleCorrection(inv *sei.ParsedInvoice, originalInvoiceData *annualregistry.Invoice) error {
	return nil
}
