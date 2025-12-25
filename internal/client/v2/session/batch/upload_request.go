package batch

import (
	"ksef/internal/client/v2/session/batch/archive"
	"ksef/internal/client/v2/session/types"
	"ksef/internal/encryption"
)

type BatchSessionPayload struct {
	Archive          *archive.Archive
	InvoiceChecksums []string
}

type batchSessionInitRequest struct {
	FormCode   types.InvoiceFormCode        `json:"formCode"`
	BatchFile  batchArchiveInfo             `json:"batchFile"`
	Encryption encryption.CipherHTTPRequest `json:"encryption"`
	Offline    bool                         `json:"offlineMode"`
	Payload    *BatchSessionPayload         `json:"-"`
}

type batchSessionPartUploadRequest struct {
	OrdinalNumber int32             `json:"ordinalNumber"`
	Method        string            `json:"method"`
	Url           string            `json:"url"`
	Headers       map[string]string `json:"headers"`
}

type batchSessionInitResponse struct {
	ReferenceNumber    string                          `json:"referenceNumber"`
	PartUploadRequests []batchSessionPartUploadRequest `json:"partUploadRequests"`
}
