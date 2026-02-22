package interactive

import (
	"context"
	"fmt"
	HTTP "ksef/internal/http"
	"net/http"
)

const endpointSessionClose = "/v2/sessions/online/%s/close"

func (s *Session) closeUploadSession(ctx context.Context, us *uploadSession) error {
	_, err := s.httpClient.Request(
		ctx,
		HTTP.RequestConfig{
			ExpectedStatus: http.StatusNoContent,
			Method:         http.MethodPost,
		},
		fmt.Sprintf(endpointSessionClose, us.refNo),
	)
	return err
}
