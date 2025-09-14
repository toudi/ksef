package auth

import (
	"context"
	HTTP "ksef/internal/http"
	"ksef/internal/sei/api/client/v2/auth/validator"
	"net/http"
)

const (
	endpointAuthChallenge = "/api/v2/auth/challenge"
)

func (t *TokenManager) getAuthChallenge(ctx context.Context) (*validator.AuthChallenge, error) {
	var authResponse validator.AuthChallenge

	_, err := t.httpClient.Request(
		ctx, HTTP.RequestConfig{
			Method:          http.MethodPost,
			ContentType:     HTTP.JSON,
			Dest:            &authResponse,
			DestContentType: HTTP.JSON,
			ExpectedStatus:  http.StatusOK,
		},
		endpointAuthChallenge,
	)

	if err != nil {
		return nil, err
	}

	return &authResponse, nil
}
