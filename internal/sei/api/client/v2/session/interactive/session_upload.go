package interactive

import (
	"context"
	HTTP "ksef/internal/http"
	"net/http"
)

type InvoiceUploadResponse struct {
	ReferenceNumber string `json:"referenceNumber"`
}

func (s *Session) uploadFile(ctx context.Context, us *uploadSession, filename string) error {
	payload, err := s.getUploadPayload(us, filename)
	if err != nil {
		return err
	}

	var resp InvoiceUploadResponse

	_, err = s.httpClient.Request(ctx, HTTP.RequestConfig{
		Headers:         s.authHeaders(),
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

	us.seiRefNumbers[filename] = resp.ReferenceNumber

	return nil
}
