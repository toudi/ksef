package invoice

import (
	"ksef/internal/money"
	"math"
)

type InvoiceItem struct {
	RowNo       int
	Before      bool
	Description string
	Unit        string
	Quantity    money.MonetaryValue
	UnitPrice   Price
	Attributes  map[string]string
}

func (ii *InvoiceItem) Amount() Amount {
	amount := Amount{}
	vatQuantizer := (1 + float64(ii.UnitPrice.Vat.Rate)/100)
	amountQuantizer := math.Pow10(ii.Quantity.DecimalPlaces + ii.UnitPrice.DecimalPlaces)

	if ii.UnitPrice.IsGross {
		// calculate amounts from gross to net
		amount.Gross = money.AmountInGrosze(float64(ii.Quantity.Amount*ii.UnitPrice.Amount) / amountQuantizer)
		amount.Net = money.AmountInGrosze((float64(amount.Gross) / 100) / vatQuantizer)
	} else {
		// calculate amounts from net to gross
		amount.Net = money.AmountInGrosze(float64(ii.Quantity.Amount*ii.UnitPrice.Amount) / amountQuantizer)
		amount.Gross = money.AmountInGrosze((float64(amount.Net) / 100) * vatQuantizer)
	}

	amount.VAT = amount.Gross - amount.Net

	return amount
}
