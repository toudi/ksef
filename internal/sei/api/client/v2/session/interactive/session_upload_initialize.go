package interactive

import (
	"context"
	"fmt"
	"ksef/internal/encryption"
	HTTP "ksef/internal/http"
	"ksef/internal/logging"
	"ksef/internal/registry"
	"net/http"
)

const (
	endpointInitializeSession    = "/api/v2/sessions/online"
	endpointSessionUploadInvoice = "/api/v2/sessions/online/%s/invoices"
)

type uploadSessionRequest struct {
	FormCode   registry.InvoiceFormCode     `json:"formCode"`
	Encryption encryption.CipherHTTPRequest `json:"encryption"`
}
type uploadSessionResponse struct {
	ReferenceNumber string `json:"referenceNumber"`
}

type uploadSession struct {
	refNo         string
	uploadUrl     string
	cipher        *encryption.Cipher
	seiRefNumbers map[string]string // mapping between local filename and KSeF reference numbers
}

func (s *Session) initialize(ctx context.Context, invoiceFormCode registry.InvoiceFormCode) (*uploadSession, error) {
	// prepare encryption data
	var err error

	cipher, err := encryption.CipherInit(32)
	if err != nil {
		return nil, err
	}

	var req uploadSessionRequest = uploadSessionRequest{
		FormCode: invoiceFormCode,
	}
	req.Encryption, err = cipher.PrepareHTTPRequestPayload(s.apiConfig.Certificate.PEM())
	if err != nil {
		return nil, err
	}
	var resp uploadSessionResponse

	_, err = s.httpClient.Request(
		ctx,
		HTTP.RequestConfig{
			Body:            req,
			ContentType:     HTTP.JSON,
			Dest:            &resp,
			DestContentType: HTTP.JSON,
			ExpectedStatus:  http.StatusCreated,
			Method:          http.MethodPost,
		},
		endpointInitializeSession,
	)

	if err != nil {
		return nil, err
	}

	return &uploadSession{
		refNo:     resp.ReferenceNumber,
		cipher:    cipher,
		uploadUrl: fmt.Sprintf(endpointSessionUploadInvoice, resp.ReferenceNumber),
	}, nil
}

func (s *Session) uploadInvoicesForForm(ctx context.Context, invoiceFormCode registry.InvoiceFormCode, files []registry.CollectionFile) error {
	us, err := s.initialize(ctx, invoiceFormCode)
	if err != nil {
		return err
	}
	for i, file := range files {
		if err = s.uploadFile(ctx, us, file.Filename); err != nil {
			logging.InteractiveLogger.Error("error uploading invoice", "counter", i, "error", err)
		}
	}
	return s.closeUploadSession(ctx, us)
}
