package invoices

import (
	"ksef/internal/client/v2/types/invoices"
	"time"
)

type SyncParams struct {
	DestPath      string
	PDF           bool
	SubjectType   invoices.SubjectType
	PageSize      int
	DateRangeType invoices.DateRangeType
	DateFrom      time.Time
}
