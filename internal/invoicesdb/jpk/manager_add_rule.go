package jpk

import (
	monthlyregistry "ksef/internal/invoicesdb/monthly-registry"
	"ksef/internal/invoicesdb/shared"
	subjectsettings "ksef/internal/invoicesdb/subject-settings"
)

func (j *JPKManager) AddItemRules(
	invoice *monthlyregistry.Invoice,
	rules []shared.JPKItemRule,
	global bool,
) (err error) {
	if global {
		// we have to persist this in the subject settings file
		if err = j.ss.Modify(func(ss *subjectsettings.SubjectSettings) error {
			for _, rule := range rules {
				ss.JPK.ItemRules = append(ss.JPK.ItemRules, subjectsettings.JPKRuleWithNIP{
					NIP:  invoice.Issuer.NIP,
					Rule: rule,
				})
			}
			return nil
		}); err != nil {
			return err
		}

		return j.ss.Save()
	}
	// let's persist the rules in the invoice registry
	if invoice.JPK == nil {
		invoice.JPK = &monthlyregistry.JPKProps{}
	}
	invoice.JPK.ItemRules = append(invoice.JPK.ItemRules, rules...)
	return j.registry.Save()
}
