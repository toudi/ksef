package status

import (
	"context"
	"fmt"
	HTTP "ksef/internal/http"
	"net/http"
)

const endpointSessionStatus = "/api/v2/sessions/%s"

func CheckSessionStatus(ctx context.Context, httpClient *HTTP.Client, sessionId string) (*StatusResponse, error) {
	var resp StatusResponse

	_, err := httpClient.Request(
		ctx,
		HTTP.RequestConfig{
			ContentType:     HTTP.JSON,
			Dest:            &resp,
			DestContentType: HTTP.JSON,
			ExpectedStatus:  http.StatusOK,
			Method:          http.MethodGet,
		},
		fmt.Sprintf(endpointSessionStatus, sessionId),
	)

	if err != nil {
		return nil, err
	}

	return nil, nil
}
