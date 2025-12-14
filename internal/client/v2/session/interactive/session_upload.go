package interactive

import (
	"context"
	sessionTypes "ksef/internal/client/v2/session/types"
	HTTP "ksef/internal/http"
	"net/http"
)

type InvoiceUploadResponse struct {
	ReferenceNumber string `json:"referenceNumber"`
}

func (s *Session) uploadFile(
	ctx context.Context,
	us *uploadSession,
	file sessionTypes.Invoice,
) error {
	payload, err := s.getUploadPayload(us, file.Filename, file.Offline)
	if err != nil {
		return err
	}

	var resp InvoiceUploadResponse

	_, err = s.httpClient.Request(ctx, HTTP.RequestConfig{
		Body:            payload,
		ContentType:     HTTP.JSON,
		Dest:            &resp,
		DestContentType: HTTP.JSON,
		ExpectedStatus:  http.StatusAccepted,
		Method:          http.MethodPost,
	}, us.uploadUrl)

	if err != nil {
		return err
	}

	us.seiRefNumbers[file.Filename] = resp.ReferenceNumber

	return nil
}
