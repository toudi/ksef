package uploader

import (
	"context"
	"errors"
	sessionTypes "ksef/internal/client/v2/session/types"
	"ksef/internal/logging"
)

var (
	errUpload                         = errors.New("upload error")
	errCannotLookupSrcInvoice         = errors.New("this really should not happen. cannot lookup source invoice in the upload queue")
	errUpdatingInvoiceReferenceNumber = errors.New("unable to update uploaded invoice reference number")
	errSavingRegistry                 = errors.New("unable to save registry")
)

func (u *Uploader) UploadInvoices(ctx context.Context) ([]*sessionTypes.UploadSessionResult, error) {
	result, err := u.uploadInvoices(ctx)
	if err != nil {
		logging.UploadLogger.Error("error uploading invoices", "err", err)
		return nil, err
	}

	return result, nil
}

// this method does the actual work of uploading invoices to KSeF
func (u *Uploader) uploadInvoices(ctx context.Context) ([]*sessionTypes.UploadSessionResult, error) {
	uploadSession, err := u.ksefClient.InteractiveSession()
	if err != nil {
		return nil, err
	}

	return uploadSession.UploadInvoices(
		ctx,
		u.Queue,
	)
}
