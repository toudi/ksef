package batch

import (
	"context"
	"errors"
	"ksef/internal/certsdb"
	"ksef/internal/client/v2/session/types"
	"ksef/internal/encryption"
	HTTP "ksef/internal/http"
	baseHTTP "net/http"
	"time"
)

const (
	endpointBatchSessionInit = "/api/v2/sessions/batch"
)

var (
	ErrCannotInitializeCiper = errors.New("unable to initialize cipher")
)

func (b *Session) UploadInvoices(
	ctx context.Context,
	payload types.UploadPayload,
) ([]*types.UploadSessionResult, error) {
	var result []*types.UploadSessionResult
	// v2 specs forces us to group invoices by their form code
	// on the other hand, it no longer forces us to send invoices through a 3rd party server
	for formCode, files := range payload {
		uploadResult, err := b.uploadInvoicesForForm(ctx, formCode, files)
		if err != nil {
			return nil, err
		}
		result = append(result, uploadResult)
	}

	return result, nil
}

type batchArchivePart struct {
	OrdNo    uint32 `json:"ordinalNumber"`
	FileName string `json:"fileName"`
	FileSize uint64 `json:"fileSize"`
	FileHash string `json:"fileHash"`
}

type batchArchiveInfo struct {
	FileSize  uint64             `json:"fileSize"`
	FileHash  string             `json:"fileHash"`
	FileParts []batchArchivePart `json:"fileParts"`
}

func (b *Session) uploadInvoicesForForm(
	ctx context.Context,
	formCode types.InvoiceFormCode,
	files []types.Invoice,
) (*types.UploadSessionResult, error) {
	// first, let's prepare metadata info
	initSessionReq, err := b.generateMetadataByFormCode(formCode, files)
	if err != nil {
		return nil, err
	}

	cipher, err := encryption.CipherInit(32)
	if err != nil {
		return nil, errors.Join(ErrCannotInitializeCiper, err)
	}

	certificate, err := b.certsDB.GetByUsage(certsdb.UsageSymmetricKeyEncryption, "")
	if err != nil {
		return nil, err
	}
	initSessionReq.Encryption, err = cipher.PrepareHTTPRequestPayload(certificate.Filename())
	if err != nil {
		return nil, err
	}

	var resp batchSessionInitResponse

	_, err = b.httpClient.Request(
		ctx,
		HTTP.RequestConfig{
			Body:           initSessionReq,
			ContentType:    HTTP.JSON,
			Dest:           resp,
			ExpectedStatus: baseHTTP.StatusCreated,
		},
		endpointBatchSessionInit,
	)

	if err != nil {
		return nil, err
	}

	// TODO: implement actual uploading
	return &types.UploadSessionResult{
		Timestamp: time.Now(),
		SessionID: resp.ReferenceNumber,
	}, nil
}
