package common

import "ksef/common/xml"

type InvoiceItem struct {
	Description    string
	Unit           string
	Quantity       int
	UnitPriceNet   int
	UnitPriceGross int
	AmountNet      int
	AmountGross    int
	VATRate        string
	Attributes     map[string]string
}

type Invoice struct {
	Number     string
	Items      []*InvoiceItem
	Issuer     *Issuer
	Attributes map[string]string
}

type Issuer struct {
	Node *xml.Node
}
