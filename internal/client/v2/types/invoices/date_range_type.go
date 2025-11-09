package invoices

type DateRangeType string

const (
	DateTypeIssue     DateRangeType = "Issue"
	DateTypeInvoicing DateRangeType = "Invoicing"
	DateTypeStorage   DateRangeType = "PermanentStorage"
)
