package invoice

import (
	"ksef/internal/money"
	"time"
)

type VAT struct {
	Rate        int
	Description string
	Except      bool
}

type Price struct {
	money.MonetaryValue
	IsGross bool
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

func (i *Invoice) Clear() {
	i.Items = make([]*InvoiceItem, 0)
	i.TotalPerVATRate = make(map[string]Amount)
	i.Attributes = make(map[string]string)
	i.Total = Amount{}
}
