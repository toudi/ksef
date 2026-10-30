package invoices

import (
	"ksef/core/ksef"
	"time"
)

type SyncParams struct {
	DestPath      string
	PDF           bool
	SubjectType   ksef.SubjectType
	PageSize      int
	DateRangeType ksef.DateRangeType
	DateFrom      time.Time
}
