package registry

import "ksef/internal/registry/types"

func (ir *InvoiceRegistry) Update(invoice types.Invoice) (types.Invoice, error) {
	index, exists := ir.checksumIndex[invoice.Checksum]

	if !exists {
		return types.Invoice{}, ErrDoesNotExist
	}

	ir.Invoices[index] = invoice
	return invoice, nil
}

func (ir *InvoiceRegistry) Upsert(invoice types.Invoice) (types.Invoice, error) {
	index, exists := ir.checksumIndex[invoice.Checksum]

	if !exists {
		ir.Invoices = append(ir.Invoices, invoice)
	} else {
		ir.Invoices[index] = invoice
	}

	return invoice, nil
}
