package interactive

import (
	"context"
	HTTP "ksef/internal/http"
	"ksef/internal/registry"
	"ksef/internal/utils"
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

	// notify invoice registry that we've uploaded `filename` within `us.refNo` upload session and that it has
	// received `resp.ReferenceNumber`. This will be important when we'll want to retrieve UPO's
	s.registry.SetUploadResult(us.refNo, &registry.InvoiceUploadResult{
		Filename: filename,
		Checksum: utils.Base64ToHex(payload.InvoiceHash),
		SeiRefNo: resp.ReferenceNumber,
	})

	return nil
}
