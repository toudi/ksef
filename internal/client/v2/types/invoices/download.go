package invoices

import (
	"time"
)

type DownloadParams struct {
	Incremental  bool
	PDF          bool
	SubjectTypes []SubjectType
	StartDate    time.Time
	EndDate      *time.Time
	PageSize     int
}
