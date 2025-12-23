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
	if index, exists := r.checksumIndex[checksum]; exists {
		return r.Invoices[index]
	}
	return nil
}

func (r *Registry) GetUnsynced() ([]*InvoiceMetadata, error) {
	var unsynced []*InvoiceMetadata

	for _, invoice := range r.Invoices {
		if invoice.Type == InvoiceTypeIssued {
			// either the invoice has KSeFReFNo (which means it was already uploaded previously)
			// or it has uploadErrors - which means it was tried to be uploaded but this
			// has failed - which means we cannot upload it again.
			if invoice.KSeFRefNo != "" || len(invoice.UploadErrors) > 0 {
				continue
			}

			if invoiceMeta, err := r.getInvoiceMetadata(invoice, invoice.OrdNum); err != nil {
				return nil, err
			} else {
				unsynced = append(unsynced, invoiceMeta)
			}

		}
	}

	return unsynced, nil
}
