package invoices

import (
	"time"
)

type DownloadParams struct {
	Incremental  bool
	PDF          bool
	SubjectTypes []SubjectType
	DateType     DateRangeType
	StartDate    time.Time
	EndDate      *time.Time
	PageSize     int
}
