package common

import "time"

type VAT struct {
	Rate        int
	Description string
}

type Price struct {
	IsGross bool
	Value   int
	Vat     VAT
}

type Amount struct {
	// these ints indicate amounts in groszy's (1/100'ths of 1PLN)
	Net   int
	Gross int
	VAT   int
}

type Invoice struct {
	Number           string
	Issued           time.Time
	Items            []*InvoiceItem
	TotalPerVATRate  map[string]Amount
	Total            Amount
	Attributes       map[string]string
	BasedOnNetPrices bool
}
