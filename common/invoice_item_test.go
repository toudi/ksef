package common

import "testing"

type Expected struct {
	Net   int
	Gross int
	Vat   int
}
type AmountCalculationTestCase struct {
	Quantity  MonetaryValue
	UnitPrice int
	IsGross   bool
	Expected  Expected
}

func TestCalculatingPricesOnInvoiceItem(t *testing.T) {
	for _, testCase := range []AmountCalculationTestCase{
		{
			Quantity: MonetaryValue{
				Amount:        2,
				DecimalPlaces: 0,
			},
			UnitPrice: 100,
			IsGross:   false,
			Expected: Expected{
				Net:   200,
				Gross: 246,
				Vat:   46,
			},
		},
		{
			Quantity: MonetaryValue{
				Amount:        2,
				DecimalPlaces: 0,
			},
			UnitPrice: 100,
			IsGross:   true,
			Expected: Expected{
				Net:   163,
				Gross: 200,
				Vat:   37,
			},
		},
		{
			Quantity: MonetaryValue{
				Amount:        200,
				DecimalPlaces: 0,
			},
			UnitPrice: 100,
			IsGross:   true,
			Expected: Expected{
				Net:   16260,
				Gross: 20000,
				Vat:   3740,
			},
		},
	} {
		item := InvoiceItem{
			Quantity: testCase.Quantity,
			UnitPrice: Price{
				MonetaryValue: MonetaryValue{
					Amount: testCase.UnitPrice,
				},
				IsGross: testCase.IsGross,
				Vat: VAT{
					Rate: 23,
				},
			},
		}

		amount := item.Amount()
		if amount.Net != testCase.Expected.Net {
			t.Errorf("unexpected net amount: %d; expected %d", amount.Net, testCase.Expected.Net)
		}
		if amount.Gross != testCase.Expected.Gross {
			t.Errorf("unexpected gross amount: %d; expected %d", amount.Gross, testCase.Expected.Gross)
		}
		if amount.VAT != testCase.Expected.Vat {
			t.Errorf("unexpected vat amount: %d; expected %d", amount.VAT, testCase.Expected.Vat)
		}
	}
}
