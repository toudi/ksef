package jpk_v7m_3

import (
	"ksef/internal/invoicesdb/jpk/types"
	"ksef/internal/money"
	"strings"
)

type FormFields struct {
	BaseField string
	VatField  string
}

type fieldToAmountRegistry map[string]money.MonetaryValue

func (far fieldToAmountRegistry) accumulate(fields []string) money.MonetaryValue {
	accumulated := money.MonetaryValue{}

	for _, field := range fields {
		if amt, exists := far[field]; exists {
			accumulated = accumulated.Add(amt)
		}
	}

	return accumulated
}

func (far fieldToAmountRegistry) Add(fields FormFields, amounts *types.VATInfo) {
	baseField := strings.Replace(fields.BaseField, "K_", "P_", 1)
	vatField := strings.Replace(fields.VatField, "K_", "P_", 1)
	if _, ok := far[baseField]; !ok {
		far[baseField] = money.MonetaryValue{}
	}
	far[baseField] = far[baseField].Add(amounts.Base)

	if vatField != "" {
		if _, ok := far[vatField]; !ok {
			far[vatField] = money.MonetaryValue{}
		}
		far[vatField] = far[vatField].Add(amounts.Vat)
	}
}
