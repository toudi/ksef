package invoices

import (
	"time"

	"ksef/internal/client/v2/types/invoices"
)

type SyncParams struct {
	DestPath      string
	PDF           bool
	SubjectType   invoices.SubjectType
	PageSize      int
	DateRangeType invoices.DateRangeType
	DateFrom      time.Time
}
