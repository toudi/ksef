package monthlyregistry

func (r *Registry) assignOrdNums() {
	for _, invoice := range r.Invoices {
		invoiceType := invoice.Type
		if _, exists := r.OrdNums[invoiceType]; !exists {
			r.OrdNums[invoiceType] = 0
		}
		r.OrdNums[invoiceType]++
		invoice.OrdNum = r.OrdNums[invoiceType]
	}
}
