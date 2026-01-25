package subjectsettings

import (
	monthlyregistry "ksef/internal/invoicesdb/monthly-registry"
	"ksef/internal/invoicesdb/shared"
)

type JPKRuleWithNIP struct {
	NIP  string             `yaml:"nip"`
	Rule shared.JPKItemRule `yaml:"rule"`
}

type JPKFormMeta struct {
	IRSCode    int               `yaml:"irs-code,omitempty"`
	SystemName string            `yaml:"system-name,omitempty"`
	Subject    map[string]any    `yaml:"subject,omitempty"`
	Defaults   map[string]string `yaml:"defaults,omitempty"`
}

type JPKSettings struct {
	ItemRules []JPKRuleWithNIP `yaml:"item-rules,omitempty"`
	FormMeta  JPKFormMeta      `yaml:"form,omitempty"`
}

func (j *JPKSettings) ItemHasRule(invoice *monthlyregistry.Invoice, itemHash shared.ItemHash, ruleChecker func(rule shared.JPKItemRule) bool) bool {
	// or it can be done via item level rules
	for _, itemRule := range j.ItemRules {
		if itemRule.NIP != invoice.Issuer.NIP {
			continue
		}

		if itemRule.Rule.Hash.Matches(itemHash) && ruleChecker(itemRule.Rule) {
			return true
		}
	}

	return false
}
