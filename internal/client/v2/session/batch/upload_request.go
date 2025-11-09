package batch

import (
	"ksef/internal/encryption"
	"ksef/internal/registry"
	"net/url"
)

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
