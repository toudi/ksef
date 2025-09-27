package batch

import (
	"context"
	"errors"
	"ksef/internal/encryption"
	HTTP "ksef/internal/http"
	"ksef/internal/registry"
	baseHTTP "net/http"
	"net/url"
)

const (
	endpointBatchSessionInit = "/api/v2/sessions/batch"
)

var (
	ErrCannotInitializeCiper = errors.New("unable to initialize cipher")
)

func (b *Session) UploadInvoices() error {
	collection, err := b.registry.InvoiceCollection()
	if err != nil {
		return err
	}

	// v2 specs forces us to group invoices by their form code
	// on the other hand, it no longer forces us to send invoices through a 3rd party server
	for formCode, files := range collection.Files {
		if err := b.uploadInvoicesForForm(context.Background(), formCode, files); err != nil {
			return err
		}
	}

	return nil

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

type batchSessionInitRequest struct {
	FormCode   registry.InvoiceFormCode     `json:"formCode"`
	BatchFile  batchArchiveInfo             `json:"batchFile"`
	Encryption encryption.CipherHTTPRequest `json:"encryption"`
}

type batchSessionPartUploadRequest struct {
	OrdinalNumber int32             `json:"ordinalNumber"`
	Method        string            `json:"method"`
	Url           url.URL           `json:"url"`
	Headers       map[string]string `json:"headers"`
}

type batchSessionInitResponse struct {
	ReferenceNumber    string                          `json:"referenceNumber"`
	PartUploadRequests []batchSessionPartUploadRequest `json:"partUploadRequests"`
}

func (b *Session) uploadInvoicesForForm(ctx context.Context, formCode registry.InvoiceFormCode, files []registry.CollectionFile) error {
	// first, let's prepare metadata info
	initSessionReq, err := b.generateMetadataByFormCode(formCode, files)
	if err != nil {
		return err
	}

	cipher, err := encryption.CipherInit(32)
	if err != nil {
		return errors.Join(ErrCannotInitializeCiper, err)
	}

	initSessionReq.Encryption, err = cipher.PrepareHTTPRequestPayload(b.apiConfig.Certificate.PEM())
	if err != nil {
		return err
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

	return err
}
