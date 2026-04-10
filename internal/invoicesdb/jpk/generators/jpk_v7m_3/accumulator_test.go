package jpk_v7m_3

import (
	"ksef/internal/invoicesdb/jpk/types"
	"ksef/internal/money"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFieldToAmountRegistryAdd(t *testing.T) {
	type sourceItem struct {
		fields FormFields
		info   types.VATInfo
	}

	type unitTest struct {
		name     string
		inputs   []sourceItem
		expected map[string]money.MonetaryValue
	}

	for _, tt := range []unitTest{
		{
			name: "single item without VAT field",
			inputs: []sourceItem{
				{fields: FormFields{BaseField: "K_10"}, info: types.VATInfo{Base: money.MonetaryValue{Amount: 1000, DecimalPlaces: 2}}},
			},
			expected: map[string]money.MonetaryValue{"P_10": {Amount: 1000, DecimalPlaces: 2}},
		},
		{
			name: "item with VAT field",
			inputs: []sourceItem{
				{fields: FormFields{BaseField: "K_15", VatField: "K_16"}, info: types.VATInfo{Base: money.MonetaryValue{Amount: 300, DecimalPlaces: 2}, Vat: money.MonetaryValue{Amount: 15, DecimalPlaces: 2}}},
			},
			expected: map[string]money.MonetaryValue{"P_15": {Amount: 300, DecimalPlaces: 2}, "P_16": {Amount: 15, DecimalPlaces: 2}},
		},
		{
			name: "two items with same fields",
			inputs: []sourceItem{
				{fields: FormFields{BaseField: "K_19", VatField: "K_20"}, info: types.VATInfo{Base: money.MonetaryValue{Amount: 100, DecimalPlaces: 2}, Vat: money.MonetaryValue{Amount: 22, DecimalPlaces: 2}}},
				{fields: FormFields{BaseField: "K_19", VatField: "K_20"}, info: types.VATInfo{Base: money.MonetaryValue{Amount: 50, DecimalPlaces: 2}, Vat: money.MonetaryValue{Amount: 11, DecimalPlaces: 2}}},
			},
			expected: map[string]money.MonetaryValue{"P_19": {Amount: 150, DecimalPlaces: 2}, "P_20": {Amount: 33, DecimalPlaces: 2}},
		},
		{
			name: "two items with same fields (7% and 8%)",
			inputs: []sourceItem{
				{fields: FormFields{BaseField: "K_17", VatField: "K_18"}, info: types.VATInfo{Base: money.MonetaryValue{Amount: 200, DecimalPlaces: 2}, Vat: money.MonetaryValue{Amount: 14, DecimalPlaces: 2}}},
				{fields: FormFields{BaseField: "K_17", VatField: "K_18"}, info: types.VATInfo{Base: money.MonetaryValue{Amount: 300, DecimalPlaces: 2}, Vat: money.MonetaryValue{Amount: 24, DecimalPlaces: 2}}},
			},
			expected: map[string]money.MonetaryValue{"P_17": {Amount: 500, DecimalPlaces: 2}, "P_18": {Amount: 38, DecimalPlaces: 2}},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			far := fieldToAmountRegistry{}

			for _, item := range tt.inputs {
				far.Add(item.fields, &item.info)
			}

			require.Equal(t, tt.expected, map[string]money.MonetaryValue(far))
		})
	}
}
