package status

import (
	"context"
	"fmt"
	ratelimits "ksef/internal/client/v2/rate-limits"
	HTTP "ksef/internal/http"
	"net/http"
)

const endpointSessionStatus = "/api/v2/sessions/%s"

func (sc *SessionStatusChecker) CheckSessionStatus(
	ctx context.Context,
	sessionId string,
) (*StatusResponse, error) {
	var resp StatusResponse

	_, err := sc.httpClient.Request(
		ctx,
		HTTP.RequestConfig{
			ContentType:     HTTP.JSON,
			Dest:            &resp,
			DestContentType: HTTP.JSON,
			ExpectedStatus:  http.StatusOK,
			Method:          http.MethodGet,
			OperationId:     ratelimits.OperationSessionMisc,
		},
		fmt.Sprintf(endpointSessionStatus, sessionId),
	)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}
