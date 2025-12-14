package batch

import (
	"ksef/internal/client/v2/session/types"
	"ksef/internal/encryption"
	"net/url"
)

type batchSessionInitRequest struct {
	FormCode   types.InvoiceFormCode        `json:"formCode"`
	BatchFile  batchArchiveInfo             `json:"batchFile"`
	Encryption encryption.CipherHTTPRequest `json:"encryption"`
	Offline    bool                         `json:"offlineMode"`
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
