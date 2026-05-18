package annotations

import (
	monthlyregistry "ksef/internal/invoicesdb/monthly-registry"
	"ksef/internal/invoicesdb/shared"
)

func (m *Annotations) ItemHasRule(invoice *monthlyregistry.Invoice, hash shared.ItemHash, rule func(rule shared.Annotation) bool) bool {
	// first check the item directly (i.e. from within the registry)
	if checkInvoiceItemHasRule(invoice, hash, rule) {
		return true
	}

	// if this did not yield result, we can still check under subject settings
	if m.ss != nil && m.ss.Annotations != nil {
		return m.ss.Annotations.ItemHasRule(invoice, hash, rule)
	}
	return false
}

func checkInvoiceItemHasRule(invoice *monthlyregistry.Invoice, hash shared.ItemHash, ruleChecker func(rule shared.Annotation) bool) bool {
	if invoice.Annotations == nil {
		return false
	}

	for _, rule := range invoice.Annotations {
		if rule.Hash.Matches(hash) && ruleChecker(rule) {
			return true
		}
	}

	return false
}

func (m *Annotations) GetItemRule(invoice *monthlyregistry.Invoice, hash shared.ItemHash) *shared.Annotation {
	// First check invoice-level rules
	if invoice.Annotations != nil {
		for _, rule := range invoice.Annotations {
			if rule.Hash.Matches(hash) {
				return &rule
			}
		}
	}

	// Then check subject-level (global) rules
	if m.ss != nil && m.ss.Annotations != nil {
		for _, itemRule := range m.ss.Annotations {
			if invoice.Issuer != nil && itemRule.NIP != invoice.Issuer.NIP {
				continue
			}
			if itemRule.Rule.Hash.Matches(hash) {
				return &itemRule.Rule
			}
		}
	}

	return nil
}
