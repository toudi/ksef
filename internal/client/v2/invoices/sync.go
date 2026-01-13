package invoices

import (
	types "ksef/internal/client/v2/types/invoices"
	"time"
)

type DateRangeType string

const (
	endpointInvoicesMetadata               = "/api/v2/invoices/query/metadata"
	DateRangeTypeIssue       DateRangeType = "Issue"
	DateRangeStorage         DateRangeType = "PermanentStorage"
)

type DateRange struct {
	DateType DateRangeType `json:"dateType"`
	From     time.Time     `json:"from"`
	To       *time.Time    `json:"to,omitzero"`
}

type InvoiceMetadataRequest struct {
	SubjectType types.SubjectType `json:"subjectType"`
	DateRange   DateRange         `json:"dateRange"`
}
