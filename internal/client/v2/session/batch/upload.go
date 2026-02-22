package batch

import (
	"context"
	"errors"
	"fmt"
	"ksef/internal/certsdb"
	"ksef/internal/client/v2/session/status"
	"ksef/internal/client/v2/session/types"
	"ksef/internal/encryption"
	HTTP "ksef/internal/http"
	"ksef/internal/logging"
	baseHTTP "net/http"
	"os"
	"time"
)

const (
	endpointBatchSessionInit  = "/v2/sessions/batch"
	endpointBatchSessionClose = "/v2/sessions/batch/%s/close"
)

var (
	ErrCannotInitializeCiper = errors.New("unable to initialize cipher")
	errEncryptionError       = errors.New("error encrypting batch part")
	errOpeningEncryptedFile  = errors.New("error opening encrypted filename")
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
	logger := logging.UploadLogger.With("session type", "batch")
	logger.Debug("initialize cipher")
	cipher, err := encryption.CipherInit(32)
	if err != nil {
		return nil, errors.Join(ErrCannotInitializeCiper, err)
	}

	// first, let's prepare metadata info
	logger.Debug("generateMetadata (generate batch session)")
	initSessionReq, err := b.generateMetadataByFormCode(formCode, files, cipher)
	if err != nil {
		return nil, err
	}

	logger.Debug("fetch symmetric key encryption certificate")
	certificate, err := b.certsDB.GetByUsage(certsdb.UsageSymmetricKeyEncryption, "")
	if err != nil {
		return nil, err
	}
	initSessionReq.Encryption, err = cipher.PrepareHTTPRequestPayload(certificate.Filename())
	if err != nil {
		return nil, err
	}

	var uploadedInvoices []status.InvoiceInfo
	logger.Debug("populate session id's for status checker")
	for _, checksum := range initSessionReq.Payload.InvoiceChecksums {
		uploadedInvoices = append(uploadedInvoices, status.InvoiceInfo{Checksum: checksum})
	}
	var resp batchSessionInitResponse

	logger.Debug("HTTP request to session init")
	_, err = b.httpClient.Request(
		ctx,
		HTTP.RequestConfig{
			Body:            initSessionReq,
			ContentType:     HTTP.JSON,
			DestContentType: HTTP.JSON,
			Method:          baseHTTP.MethodPost,
			Dest:            &resp,
			ExpectedStatus:  baseHTTP.StatusCreated,
		},
		endpointBatchSessionInit,
	)
	if err != nil {
		return nil, err
	}

	logger.Debug("uploading parts")
	unauthedClient := HTTP.NewClient("")
	for partIdx, reqParams := range resp.PartUploadRequests {
		batchPart := initSessionReq.Payload.Archive.Parts[partIdx]
		encryptedReader, err := os.Open(batchPart.EncryptedFilename)
		if err != nil {
			return nil, errors.Join(errOpeningEncryptedFile, err)
		}
		logger.Debug("upload", "batchPart", batchPart)
		_, err = unauthedClient.Request(
			ctx,
			HTTP.RequestConfig{
				Method:         reqParams.Method,
				Headers:        reqParams.Headers,
				ContentType:    HTTP.BIN,
				Body:           encryptedReader,
				ExpectedStatus: baseHTTP.StatusCreated,
			},
			reqParams.Url,
		)
		encryptedReader.Close()
		if err != nil {
			return nil, err
		}
	}

	logger.Debug("HTTP close batch session")
	if _, err = b.httpClient.Request(
		ctx,
		HTTP.RequestConfig{
			Method:         baseHTTP.MethodPost,
			ExpectedStatus: baseHTTP.StatusNoContent,
		},
		fmt.Sprintf(endpointBatchSessionClose, resp.ReferenceNumber),
	); err != nil {
		return nil, err
	}

	logger.Debug("cleanup temporary dir")
	if err = initSessionReq.Payload.Archive.Cleanup(); err != nil {
		return nil, err
	}

	result := &types.UploadSessionResult{
		Timestamp: time.Now(),
		SessionID: resp.ReferenceNumber,
		Invoices:  uploadedInvoices,
	}
	logger.Debug("all fine - return result", "result", result)
	return result, nil
}
