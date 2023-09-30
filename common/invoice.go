package common

type VAT struct {
	Exempt bool
	Rate   int
}

type Price struct {
	IsGross bool
	Value   int
	Vat     VAT
}

type InvoiceItem struct {
	Description string
	Unit        string
	Quantity    float64
	UnitPrice   Price
	Attributes  map[string]string
}

type Invoice struct {
	Number     string
	Items      []*InvoiceItem
	Attributes map[string]string
}
