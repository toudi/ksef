package typst

import (
	"errors"
	"fmt"
	"ksef/internal/config/pdf/abstract"
	"text/template"
	"time"

	"github.com/spf13/viper"
)

var funcs = template.FuncMap{
	"add": func(value int) int {
		return value + 1
	},
	"float2string": func(value float64) string {
		return fmt.Sprintf("%.0f", value)
	},
	"grosz2pln": func(value int) string {
		return fmt.Sprintf("%.2f", float64(value)/100)
	},
	"map": func(pairs ...any) (map[string]any, error) {
		if len(pairs)%2 != 0 {
			return nil, errors.New("misaligned map")
		}

		m := make(map[string]any, len(pairs)/2)

		for i := 0; i < len(pairs); i += 2 {
			key, ok := pairs[i].(string)

			if !ok {
				return nil, fmt.Errorf("cannot use type %T as map key", pairs[i])
			}
			m[key] = pairs[i+1]
		}
		return m, nil
	},
	"has": func(a_map map[string]any, key string) bool {
		_, exists := a_map[key]
		return exists
	},
	"datetime_hr": func(timestamp time.Time) string {
		return timestamp.Format(`2006-01-02 \ 15:04:05 MST`)
	},
}

func parseTemplates(vip *viper.Viper) (templates abstract.Templates, err error) {
	invoiceTemplatesPath := vip.GetString(cfgKeyTypstInvoiceTemplates)
	upoTemplatesPath := vip.GetString(cfgKeyTypstUpoTemplates)

	var isValid = false

	if invoiceTemplatesPath != "" {
		templates.Invoice, err = abstract.ReadTemplatesFromDirectory(invoiceTemplatesPath, "typ", funcs)
		isValid = err == nil

		templates.Invoice.Header.Left = vip.GetString(cfgKeyTypstInvoiceHeaderL)
		templates.Invoice.Header.Center = vip.GetString(cfgKeyTypstInvoiceHeaderC)
		templates.Invoice.Header.Right = vip.GetString(cfgKeyTypstInvoiceHeaderR)

		templates.Invoice.Footer.Left = vip.GetString(cfgKeyTypstInvoiceFooterL)
		templates.Invoice.Footer.Center = vip.GetString(cfgKeyTypstInvoiceFooterC)
		templates.Invoice.Footer.Right = vip.GetString(cfgKeyTypstInvoiceFooterR)
	}

	if upoTemplatesPath != "" {
		templates.UPO, err = abstract.ReadTemplatesFromDirectory(upoTemplatesPath, "typ", funcs)
		isValid = err == nil
	}

	if err == nil && !isValid {
		err = abstract.ErrNoTemplatesDefined
	}

	return templates, err
}
