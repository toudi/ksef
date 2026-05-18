package monthlyregistry

// JPKEligibleInvoices returns invoices that are eligible for JPK export.
// An invoice is eligible if:
//   - its type is Issued or Received
//   - it is not completely excluded from JPK processing
func (r *Registry) JPKEligibleInvoices() []*Invoice {
	var result []*Invoice
	for _, invoice := range r.Invoices {
		if !isInvoiceJPKEligible(invoice) {
			continue
		}
		result = append(result, invoice)
	}
	return result
}

// isInvoiceJPKEligible checks if a single invoice is eligible for JPK export.
func isInvoiceJPKEligible(invoice *Invoice) bool {
	if invoice.Type != InvoiceTypeIssued && invoice.Type != InvoiceTypeReceived {
		return false
	}
	if isInvoiceExcluded(invoice.Annotations) {
		return false
	}
	return true
}

// isInvoiceExcluded checks if an invoice is completely excluded from JPK
// processing (i.e. has a single wildcard rule with Exclude=true).
func isInvoiceExcluded(annotations Annotations) bool {
	if annotations == nil || len(annotations) != 1 {
		return false
	}
	rule := annotations[0]
	return rule.Hash.Wildcard && rule.Exclude
}
