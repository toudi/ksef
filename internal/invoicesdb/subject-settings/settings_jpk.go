package subjectsettings

import (
	monthlyregistry "ksef/internal/invoicesdb/monthly-registry"
	"ksef/internal/invoicesdb/shared"
)

type JPKRuleWithNIP struct {
	NIP  string             `yaml:"nip"`
	Rule shared.JPKItemRule `yaml:"rule"`
}

type JPKSettings struct {
	ItemRules []JPKRuleWithNIP `yaml:"item-rules,omitempty"`
}

func (j *JPKSettings) IsItemExcluded(invoice *monthlyregistry.Invoice, itemHash shared.ItemHash) bool {
	// or it can be done via item level rules
	for _, itemRule := range j.ItemRules {
		if itemRule.NIP != invoice.Issuer.NIP {
			continue
		}

		if itemRule.Rule.Hash.Matches(itemHash) && itemRule.Rule.Exclude {
			return true
		}
	}

	return false
}

func (j *JPKSettings) IsItemVat50Percent(invoice *monthlyregistry.Invoice, itemHash shared.ItemHash) bool {
	// or it can be done via item level rules
	for _, itemRule := range j.ItemRules {
		if itemRule.NIP != invoice.Issuer.NIP {
			continue
		}

		if itemRule.Rule.Hash.Matches(itemHash) && itemRule.Rule.Vat50Percent {
			return true
		}
	}

	return false
}
