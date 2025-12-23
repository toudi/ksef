package monthlyregistry

func (r *Registry) InvoiceFilename(i *Invoice) string {
	return r.getIssuedInvoiceFilename(
		i.RefNo,
		i.OrdNum,
	)
}
