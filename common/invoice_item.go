package common

type InvoiceItem struct {
	Description string
	Unit        string
	Quantity    float64
	UnitPrice   Price
	Attributes  map[string]string
}

func (ii *InvoiceItem) Amount() Amount {
	amount := Amount{}
	vatQuantizer := (1 + float64(ii.UnitPrice.Vat.Rate)/100)

	if ii.UnitPrice.IsGross {
		// calculate amounts from gross to net
		amount.Gross = AmountInGrosze(ii.Quantity * float64(ii.UnitPrice.Value) / 100)
		amount.Net = AmountInGrosze((float64(amount.Gross) / 100) / vatQuantizer)
	} else {
		// calculate amounts from net to gross
		amount.Net = AmountInGrosze(ii.Quantity * float64(ii.UnitPrice.Value) / 100)
		amount.Gross = AmountInGrosze((float64(amount.Net) / 100) * vatQuantizer)
	}

	amount.VAT = amount.Gross - amount.Net

	return amount
}
