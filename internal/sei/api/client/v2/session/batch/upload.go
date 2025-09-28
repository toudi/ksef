package batch

import (
	"context"
	"errors"
	"ksef/internal/certsdb"
	"ksef/internal/encryption"
	HTTP "ksef/internal/http"
	"ksef/internal/registry"
	baseHTTP "net/http"
)

const (
	endpointBatchSessionInit = "/api/v2/sessions/batch"
)

var (
	ErrCannotInitializeCiper = errors.New("unable to initialize cipher")
)

func (b *Session) UploadInvoices(ctx context.Context) error {
	collection, err := b.registry.InvoiceCollection()
	if err != nil {
		return err
	}

	// v2 specs forces us to group invoices by their form code
	// on the other hand, it no longer forces us to send invoices through a 3rd party server
	for formCode, files := range collection.Files {
		if err := b.uploadInvoicesForForm(ctx, formCode, files); err != nil {
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

	certificate, err := b.apiConfig.CertificatesDB.GetByUsage(certsdb.UsageSymmetricKeyEncryption)
	if err != nil {
		return err
	}
	initSessionReq.Encryption, err = cipher.PrepareHTTPRequestPayload(certificate.PEMFile)
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
