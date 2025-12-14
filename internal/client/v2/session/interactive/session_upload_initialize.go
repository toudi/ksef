package interactive

import (
	"context"
	"fmt"
	"ksef/internal/certsdb"
	"ksef/internal/client/v2/session/types"
	"ksef/internal/encryption"
	HTTP "ksef/internal/http"
	"ksef/internal/logging"
	"net/http"
)

const (
	endpointInitializeSession    = "/api/v2/sessions/online"
	endpointSessionUploadInvoice = "/api/v2/sessions/online/%s/invoices"
)

type uploadSessionRequest struct {
	FormCode   types.InvoiceFormCode        `json:"formCode"`
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

func (s *Session) initialize(ctx context.Context, invoiceFormCode types.InvoiceFormCode) (*uploadSession, error) {
	// prepare encryption data
	var err error

	cipher, err := encryption.CipherInit(32)
	if err != nil {
		return nil, err
	}

	var req uploadSessionRequest = uploadSessionRequest{
		FormCode: invoiceFormCode,
	}
	certificate, err := s.certsDB.GetByUsage(certsdb.UsageSymmetricKeyEncryption, "")
	if err != nil {
		return nil, err
	}
	req.Encryption, err = cipher.PrepareHTTPRequestPayload(certificate.Filename())
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
		refNo:         resp.ReferenceNumber,
		cipher:        cipher,
		uploadUrl:     fmt.Sprintf(endpointSessionUploadInvoice, resp.ReferenceNumber),
		seiRefNumbers: make(map[string]string),
	}, nil
}

func (s *Session) uploadInvoicesForForm(
	ctx context.Context,
	invoiceFormCode types.InvoiceFormCode,
	files []types.Invoice,
) (*types.UploadSessionResult, error) {
	us, err := s.initialize(ctx, invoiceFormCode)
	if err != nil {
		return nil, err
	}
	for i, file := range files {
		if err = s.uploadFile(ctx, us, file); err != nil {
			logging.InteractiveLogger.Error("error uploading invoice", "counter", i, "error", err)
		}
	}
	if err = s.closeUploadSession(ctx, us); err != nil {
		return nil, err
	}
	// collect reference numbers and/or potential errors
	var result = &types.UploadSessionResult{
		SessionID: us.refNo,
		Invoices:  make([]types.InvoiceUploadResult, 0, len(files)),
	}
	for _, file := range files {
		result.Invoices = append(result.Invoices, types.InvoiceUploadResult{
			Checksum:  file.Checksum,
			KSeFRefNo: us.seiRefNumbers[file.Filename],
		})
	}

	return result, nil
}
