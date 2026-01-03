package jpk

import (
	monthlyregistry "ksef/internal/invoicesdb/monthly-registry"
	"ksef/internal/invoicesdb/shared"
)

func (m *JPKManager) ItemIsExcluded(invoice *monthlyregistry.Invoice, hash shared.ItemHash) bool {
	// let's check the mark on the level of invoice itself
	if checkInvoiceItemExcluded(invoice, hash) {
		// if the item is marked on the invoice level - that's enough. we don't have
		// to check any further
		return true
	}

	// let's inspect subject settings
	if m.ss.JPK.IsItemExcluded(invoice, hash) {
		return true
	}
	return false
}

func (m *JPKManager) ItemHasVat50PercentFlag(invoice *monthlyregistry.Invoice, hash shared.ItemHash) bool {
	// let's check the mark on the level of invoice itself
	if checkInvoiceItemVat50Percent(invoice, hash) {
		// if the item is marked on the invoice level - that's enough. we don't have
		// to check any further
		return true
	}

	// let's inspect subject settings
	if m.ss.JPK.IsItemVat50Percent(invoice, hash) {
		return true
	}
	return false
}

func checkInvoiceItemExcluded(invoice *monthlyregistry.Invoice, hash shared.ItemHash) bool {
	if invoice.JPK == nil {
		return false
	}

	for _, rule := range invoice.JPK.ItemRules {
		if rule.Hash.Matches(hash) && rule.Exclude {
			return true
		}
	}

	return false
}

func checkInvoiceItemVat50Percent(invoice *monthlyregistry.Invoice, hash shared.ItemHash) bool {
	if invoice.JPK == nil {
		return false
	}

	for _, rule := range invoice.JPK.ItemRules {
		if rule.Hash.Matches(hash) && rule.Vat50Percent {
			return true
		}
	}

	return false
}
