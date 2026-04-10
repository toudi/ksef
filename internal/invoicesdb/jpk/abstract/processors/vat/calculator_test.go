package vat

import (
	"ksef/internal/money"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
)

func TestCalculateVat(t *testing.T) {
	type unitTest struct {
		label    string
		input    *InvoiceItem
		expected *VatInfo
		err      error
	}
	t.Run("Net -> Gross", func(t *testing.T) {
		for _, test := range []unitTest{
			{
				label: "100 NET + 23 VAT => 123 GROSS",
				input: &InvoiceItem{NetAmount: lo.ToPtr("100"), TaxRate: "23"},
				expected: &VatInfo{
					NetAmount: money.MonetaryValue{Amount: 10000, DecimalPlaces: 2},
					VatAmount: money.MonetaryValue{Amount: 2300, DecimalPlaces: 2},
					VatRate:   "23",
				},
			},
			{
				label: "81.30 NET + 18.70 VAT => 100 GROSS",
				input: &InvoiceItem{NetAmount: lo.ToPtr("81.30"), TaxRate: "23"},
				expected: &VatInfo{
					NetAmount: money.MonetaryValue{Amount: 8130, DecimalPlaces: 2},
					VatAmount: money.MonetaryValue{Amount: 1870, DecimalPlaces: 2},
					VatRate:   "23",
				},
			},
			{
				label: "100 NET + 8 VAT => 108 GROSS",
				input: &InvoiceItem{NetAmount: lo.ToPtr("100"), TaxRate: "8"},
				expected: &VatInfo{
					NetAmount: money.MonetaryValue{Amount: 10000, DecimalPlaces: 2},
					VatAmount: money.MonetaryValue{Amount: 800, DecimalPlaces: 2},
					VatRate:   "8",
				},
			},
			{
				label: "22 NET + NP VAT => 22 GROSS",
				input: &InvoiceItem{NetAmount: lo.ToPtr("22"), TaxRate: "NP"},
				expected: &VatInfo{
					NetAmount: money.MonetaryValue{Amount: 2200, DecimalPlaces: 2},
					VatRate:   "NP",
				},
			},
		} {
			t.Run(test.label, func(t *testing.T) {
				t.Parallel()

				vatInfo, err := calculateVat(test.input)
				require.Equal(t, test.expected, vatInfo)
				require.Equal(t, test.err, err)
			})
		}
	})

	t.Run("Gross -> Net", func(t *testing.T) {
		for _, test := range []unitTest{
			{
				label: "123 GROSS = 100 NET (23 VAT)",
				input: &InvoiceItem{GrossAmount: lo.ToPtr("123"), TaxRate: "23"},
				expected: &VatInfo{
					NetAmount: money.MonetaryValue{Amount: 10000, DecimalPlaces: 2},
					VatAmount: money.MonetaryValue{Amount: 2300, DecimalPlaces: 2},
					VatRate:   "23",
				},
			},
			{
				label: "123.45 GROSS = 100.37 NET (23 VAT)",
				input: &InvoiceItem{GrossAmount: lo.ToPtr("123.45"), TaxRate: "23"},
				expected: &VatInfo{
					NetAmount: money.MonetaryValue{Amount: 10037, DecimalPlaces: 2},
					VatAmount: money.MonetaryValue{Amount: 2308, DecimalPlaces: 2},
					VatRate:   "23",
				},
			},
			{
				label: "100 GROSS = 81.30 NET (23 VAT)",
				input: &InvoiceItem{GrossAmount: lo.ToPtr("100"), TaxRate: "23"},
				expected: &VatInfo{
					NetAmount: money.MonetaryValue{Amount: 8130, DecimalPlaces: 2},
					VatAmount: money.MonetaryValue{Amount: 1870, DecimalPlaces: 2},
					VatRate:   "23",
				},
			},
			{
				label: "108 GROSS = 100 NET (8 VAT)",
				input: &InvoiceItem{GrossAmount: lo.ToPtr("108"), TaxRate: "8"},
				expected: &VatInfo{
					NetAmount: money.MonetaryValue{Amount: 10000, DecimalPlaces: 2},
					VatAmount: money.MonetaryValue{Amount: 800, DecimalPlaces: 2},
					VatRate:   "8",
				},
			},
			{
				label: "22 GROSS = 22 NET (NP VAT)",
				input: &InvoiceItem{GrossAmount: lo.ToPtr("22"), TaxRate: "NP"},
				expected: &VatInfo{
					NetAmount: money.MonetaryValue{Amount: 2200, DecimalPlaces: 2},
					VatRate:   "NP",
				},
			},
		} {
			t.Run(test.label, func(t *testing.T) {
				t.Parallel()

				vatInfo, err := calculateVat(test.input)
				require.Equal(t, test.expected, vatInfo)
				require.Equal(t, test.err, err)
			})
		}
	})

	t.Run("invalid inputs", func(t *testing.T) {
		for _, test := range []unitTest{
			{
				label: "no amounts provided",
				input: &InvoiceItem{TaxRate: "23"},
				err:   ErrAmountNotDefined,
			},
		} {
			t.Run(test.label, func(t *testing.T) {
				t.Parallel()

				vatInfo, err := calculateVat(test.input)
				require.Nil(t, vatInfo)
				require.ErrorIs(t, err, test.err)
			})
		}
	})
}
