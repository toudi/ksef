package registry

func (ir *InvoiceRegistry) Update(invoice Invoice) (Invoice, error) {
	index, exists := ir.checksumIndex[invoice.Checksum]

	if !exists {
		return Invoice{}, ErrDoesNotExist
	}

	ir.Invoices[index] = invoice
	return invoice, nil
}

func (ir *InvoiceRegistry) Upsert(invoice Invoice) (Invoice, error) {
	index, exists := ir.checksumIndex[invoice.Checksum]

	if !exists {
		ir.Invoices = append(ir.Invoices, invoice)
	} else {
		ir.Invoices[index] = invoice
	}

	return invoice, nil
}
