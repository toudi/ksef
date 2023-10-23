package fa_2

// describe state that the generator is in
const (
	stateUnknown          = -1
	stateParseHeader      = iota
	stateParseInvoice     = iota
	stateParseInvoiceItem = iota
)
