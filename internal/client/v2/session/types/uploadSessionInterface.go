package types

import (
	"context"
	"ksef/internal/client/v2/session/status"
	"time"
)

type UploadSessionResult struct {
	Timestamp time.Time
	SessionID string
	Invoices  []status.InvoiceInfo
	Status    *status.StatusResponse
}

func (usr UploadSessionResult) IsProcessed() bool {
	if usr.Status == nil {
		return false
	}

	return usr.Status.IsProcessed()
}

type UploadSession interface {
	UploadInvoices(ctx context.Context, payload UploadPayload) ([]*UploadSessionResult, error)
}
