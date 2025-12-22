package monthlyregistry

func (r *Registry) ContainsHash(checksum string) bool {
	return r.GetInvoiceByChecksum(checksum) != nil
}

func (r *Registry) getInvoiceByRefNo(refNo string) *Invoice {
	// TODO: implement index. for now - just don't bother.
	for _, invoice := range r.Invoices {
		if invoice.RefNo == refNo {
			return invoice
		}
	}

	return nil
}

func (r *Registry) GetInvoiceByChecksum(checksum string) *Invoice {
	// TODO: implement proper indexing ..
	for _, invoice := range r.Invoices {
		if invoice.Checksum == checksum {
			return invoice
		}
	}

	return nil
}

func (r *Registry) GetUnsynced() ([]*InvoiceMetadata, error) {
	var issuedOrdNo int
	var unsynced []*InvoiceMetadata

	for _, invoice := range r.Invoices {
		if invoice.Type == InvoiceTypeIssued {
			issuedOrdNo += 1

			if invoice.KSeFRefNo != "" {
				continue
			}

			if invoiceMeta, err := r.getInvoiceMetadata(invoice, issuedOrdNo); err != nil {
				return nil, err
			} else {
				unsynced = append(unsynced, invoiceMeta)
			}

		}
	}

	return unsynced, nil
}
