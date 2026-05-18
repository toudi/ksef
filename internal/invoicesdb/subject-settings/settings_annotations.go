package subjectsettings

import (
	monthlyregistry "ksef/internal/invoicesdb/monthly-registry"
	"ksef/internal/invoicesdb/shared"
)

type AnnotationRuleWithNIP struct {
	NIP  string            `yaml:"nip"`
	Rule shared.Annotation `yaml:"rule"`
}

// GlobalAnnotations is a slice of annotation rules that are not tied to any particular invoice.
// instead, they describe generic rules so that if you receive invoices on a monthly basis
// these rules can be inspected upon to check whether they could be applied on the new invoices.
type GlobalAnnotations []AnnotationRuleWithNIP

func (a GlobalAnnotations) MarshalYAML() (any, error) {
	if len(a) == 0 {
		return nil, nil
	}
	return []AnnotationRuleWithNIP(a), nil
}

func (a *GlobalAnnotations) UnmarshalYAML(unmarshal func(any) error) error {
	var rules []AnnotationRuleWithNIP
	if err := unmarshal(&rules); err != nil {
		return err
	}
	*a = rules
	return nil
}

func (a GlobalAnnotations) ItemHasRule(invoice *monthlyregistry.Invoice, itemHash shared.ItemHash, ruleChecker func(rule shared.Annotation) bool) bool {
	for _, itemRule := range a {
		if itemRule.NIP != invoice.Issuer.NIP {
			continue
		}

		if itemRule.Rule.Hash.Matches(itemHash) && ruleChecker(itemRule.Rule) {
			return true
		}
	}

	return false
}
