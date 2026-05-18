package annotations

import (
	monthlyregistry "ksef/internal/invoicesdb/monthly-registry"
	"ksef/internal/invoicesdb/shared"
	subjectsettings "ksef/internal/invoicesdb/subject-settings"
)

func (a *Annotations) AddItemRules(
	invoice *monthlyregistry.Invoice,
	rules []shared.Annotation,
	global bool,
) (err error) {
	if global {
		// we have to persist this in the subject settings file
		if err = a.ss.Modify(func(ss *subjectsettings.SubjectSettings) error {
			for _, rule := range rules {
				ss.Annotations = append(ss.Annotations, subjectsettings.AnnotationRuleWithNIP{
					NIP:  invoice.Issuer.NIP,
					Rule: rule,
				})
			}
			return nil
		}); err != nil {
			return err
		}

		return a.ss.Save()
	}
	// let's persist the rules in the invoice registry
	if invoice.Annotations == nil {
		invoice.Annotations = monthlyregistry.Annotations{}
	}
	invoice.Annotations = append(invoice.Annotations, rules...)
	return a.registry.Save()
}
