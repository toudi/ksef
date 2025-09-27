package invoices

import "ksef/internal/sei/api/client/v2/types/invoices"

type SyncParams struct {
	DestPath    string
	PDF         bool
	SubjectType invoices.SubjectType
	PageSize    int
}
