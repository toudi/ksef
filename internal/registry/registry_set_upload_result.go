package registry

import (
	"errors"
	"ksef/internal/registry/types"
)

var (
	ErrUnknownUploadSessionId = errors.New("unknown upload session ID")
)

// called during upload session after successful upload, but before session is closed
func (r *InvoiceRegistry) SetUploadResult(
	uploadSessionRefNo string,
	uploadResult *types.InvoiceUploadResult,
) {
	if r.UploadSessions == nil {
		r.UploadSessions = make(map[string]*types.UploadSessionStatus)
	}

	if _, exists := r.UploadSessions[uploadSessionRefNo]; !exists {
		r.UploadSessions[uploadSessionRefNo] = &types.UploadSessionStatus{}
	}

	r.UploadSessions[uploadSessionRefNo].Invoices = append(r.UploadSessions[uploadSessionRefNo].Invoices, uploadResult)

	// TODO: save file to a temp directory to prevent data loss ?
}

func (r *InvoiceRegistry) MarkFailedInvoices(
	uploadSessionRefNo string,
	failedInvoices []int,
) error {
	sessionUploadStatus, exists := r.UploadSessions[uploadSessionRefNo]
	if !exists {
		return ErrUnknownUploadSessionId
	}

	for _, uploadResultIndex := range failedInvoices {
		sessionUploadStatus.Invoices[uploadResultIndex].Failed = true
	}

	return nil
}

// called after session is successfully processed. we can then upsert invoices
func (r *InvoiceRegistry) MarkUploadSessionProcessed(uploadSessionRefNo string) error {
	for _, uploadResult := range r.UploadSessions[uploadSessionRefNo].Invoices {
		// if invoice wasn't processed successfuly we shouldn't assign seiRefNumber to it
		if uploadResult.Failed {
			continue
		}

		invoice, err := r.GetInvoiceByChecksum(uploadResult.Checksum)
		if err != nil {
			return err
		}
		invoice.KSeFReferenceNumber = uploadResult.SeiRefNo
		if _, err = r.Update(invoice); err != nil {
			return err
		}
	}

	r.UploadSessions[uploadSessionRefNo].Processed = true

	return nil
}
