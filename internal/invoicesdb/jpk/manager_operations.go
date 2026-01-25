package jpk

import (
	monthlyregistry "ksef/internal/invoicesdb/monthly-registry"
	"ksef/internal/invoicesdb/shared"
)

func (m *JPKManager) ItemHasRule(invoice *monthlyregistry.Invoice, hash shared.ItemHash, rule func(shared.JPKItemRule) bool) bool {
	// first check the item directly (i.e. from within the registry)
	if checkInvoiceItemHasRule(invoice, hash, rule) {
		return true
	}

	// if this did not yield result, we can still check under subject settings
	return m.ss.JPK.ItemHasRule(invoice, hash, rule)
}

func checkInvoiceItemHasRule(invoice *monthlyregistry.Invoice, hash shared.ItemHash, ruleChecker func(shared.JPKItemRule) bool) bool {
	if invoice.JPK == nil {
		return false
	}

	for _, rule := range invoice.JPK.ItemRules {
		if rule.Hash.Matches(hash) && ruleChecker(rule) {
			return true
		}
	}

	return false
}
