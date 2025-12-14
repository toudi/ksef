package types

import (
	"context"
	"ksef/internal/client/v2/session/status"
)

type InvoiceUploadResult struct {
	Checksum  string
	KSeFRefNo string
	Errors    []string
}

type UploadSessionResult struct {
	SessionID string
	Invoices  []InvoiceUploadResult
	Status    *status.StatusResponse
	Processed bool
}

type UploadSession interface {
	UploadInvoices(ctx context.Context, payload UploadPayload) ([]*UploadSessionResult, error)
}
