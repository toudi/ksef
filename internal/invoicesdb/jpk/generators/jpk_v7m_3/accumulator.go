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

var vatRateToFields = map[types.VatRate]FormFields{
	types.VatRateZw:      {BaseField: "P_10"},
	types.VatRateNpI:     {BaseField: "P_11"},
	types.VatRateNpII:    {BaseField: "P_12"},
	types.VatRateZeroKR:  {BaseField: "P_13"},
	types.VatRateZeroWDT: {BaseField: "P_21"},
	types.VatRateZeroEX:  {BaseField: "P_22"},
	types.VatRate5:       {BaseField: "P_15", VatField: "P_16"},
	types.VatRate7:       {BaseField: "P_17", VatField: "P_18"},
	types.VatRate8:       {BaseField: "P_17", VatField: "P_18"},
	types.VatRate22:      {BaseField: "P_19", VatField: "P_20"},
	types.VatRate23:      {BaseField: "P_19", VatField: "P_20"},
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
