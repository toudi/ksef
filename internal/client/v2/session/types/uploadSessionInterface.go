package types

import (
	"context"
	"errors"
	"ksef/internal/client/v2/session/status"
	"ksef/internal/utils"
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

func (usr UploadSessionResult) GetInvoiceByChecksum(checksumHex string) (status.InvoiceInfo, error) {
	checksumBase64 := utils.HexToBase64(checksumHex)
	for _, invoice := range usr.Invoices {
		if invoice.ChecksumBase64 == checksumBase64 {
			return invoice, nil
		}
	}

	return status.InvoiceInfo{}, errors.New("unable to find invoice by base64 checksum")
}

type UploadSession interface {
	UploadInvoices(ctx context.Context, payload UploadPayload) ([]*UploadSessionResult, error)
}
