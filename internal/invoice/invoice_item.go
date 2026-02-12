package invoice

import (
	"ksef/internal/money"
	"math"
	"reflect"
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

func (ii InvoiceItem) IsEmpty() bool {
	// this may be tricky, however we've got a (sort of) nasty trick up uour sleves
	// what we can do is we can initialize an empty item, set the row number to the one
	// we're comparing and compare that instead. so the end result would be to compare
	// an item that consists of all of the zero values and checking if this one (ii)
	// also consists of all of the zero values (apart from the row number)
	emptyItem := InvoiceItem{RowNo: ii.RowNo, Before: ii.Before}

	return reflect.DeepEqual(ii, emptyItem)
}
