package annualregistry

func (r *Registry) GetByRefNo(refNo string) *Invoice {
	for _, inv := range r.invoices {
		if inv.RefNo == refNo {
			return inv
		}

		for _, correction := range inv.Corrections {
			if correction.RefNo == refNo {
				return correction
			}
		}

	}
	return nil
}

func (r *Registry) GetByChecksum(checksum string) *Invoice {
	for _, inv := range r.invoices {
		if inv.Checksum == checksum {
			return inv
		}

		for _, correction := range inv.Corrections {
			if correction.Checksum == checksum {
				return correction
			}
		}
	}

	return nil
}
