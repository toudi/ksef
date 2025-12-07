package annualregistry

func (r *Registry) GetByRefNo(refNo string) *Invoice {
	for _, inv := range r.invoices {
		if inv.RefNo == refNo {
			return inv
		}
	}
	return nil
}

func (r *Registry) GetByChecksum(checksum string) *Invoice {
	for _, inv := range r.invoices {
		if inv.Checksum == checksum {
			return inv
		}
	}

	return nil
}
