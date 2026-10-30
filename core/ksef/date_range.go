package ksef

type DateRangeType string

const (
	DateTypeIssue     DateRangeType = "Issue"
	DateTypeInvoicing DateRangeType = "Invoicing"
	DateTypeStorage   DateRangeType = "PermanentStorage"
)
