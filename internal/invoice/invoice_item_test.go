package invoice

import (
	"ksef/internal/money"
	"testing"
)

type Expected struct {
	Net   int
	Gross int
	Vat   int
}
type AmountCalculationTestCase struct {
	Quantity  money.MonetaryValue
	UnitPrice int
	IsGross   bool
	Expected  Expected
}

func TestCalculatingPricesOnInvoiceItem(t *testing.T) {
	for idx, testCase := range []AmountCalculationTestCase{
		{
			Quantity: money.MonetaryValue{
				Amount:        2,
				DecimalPlaces: 0,
			},
			UnitPrice: 100,
			IsGross:   false,
			Expected: Expected{
				Net:   20000,
				Gross: 24600,
				Vat:   4600,
			},
		},
		{
			Quantity: money.MonetaryValue{
				Amount:        2,
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
		{
			Quantity: money.MonetaryValue{
				Amount:        200,
				DecimalPlaces: 0,
			},
			UnitPrice: 100,
			IsGross:   true,
			Expected: Expected{
				Net:   1626016,
				Gross: 2000000,
				Vat:   373984,
			},
		},
	} {
		item := InvoiceItem{
			Quantity: testCase.Quantity,
			UnitPrice: Price{
				MonetaryValue: money.MonetaryValue{
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
			t.Fatalf("unexpected net amount at %d: %d; expected %d", idx, amount.Net, testCase.Expected.Net)
		}
		if amount.Gross != testCase.Expected.Gross {
			t.Fatalf("unexpected gross amount at %d: %d; expected %d", idx, amount.Gross, testCase.Expected.Gross)
		}
		if amount.VAT != testCase.Expected.Vat {
			t.Fatalf("unexpected vat amount at %d: %d; expected %d", idx, amount.VAT, testCase.Expected.Vat)
		}
	}
}
