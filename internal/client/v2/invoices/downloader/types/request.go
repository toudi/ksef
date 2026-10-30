package downloadertypes

import (
	"ksef/core/ksef"
	"time"
)

type DateRangeType string

const (
	DateRangeTypeIssue DateRangeType = "Issue"
	DateRangeStorage   DateRangeType = "PermanentStorage"
)

type DateRange struct {
	DateType DateRangeType `json:"dateType"`
	From     time.Time     `json:"from"`
	To       *time.Time    `json:"to,omitzero"`
}

type InvoiceListRequest struct {
	SubjectType ksef.SubjectType `json:"subjectType"`
	DateRange   DateRange        `json:"dateRange"`
}
