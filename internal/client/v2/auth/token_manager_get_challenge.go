package auth

import (
	"context"
	"ksef/internal/client/v2/auth/validator"
	HTTP "ksef/internal/http"
	"ksef/internal/logging"
	"net/http"
)

const (
	endpointAuthChallenge = "/api/v2/auth/challenge"
)

func (t *TokenManager) getAuthChallenge(ctx context.Context) (*validator.AuthChallenge, error) {
	logging.AuthLogger.Debug("get auth challenge")

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

	logging.AuthLogger.Debug("received response", "response", authResponse)

	if err != nil {
		return nil, err
	}

	return &authResponse, nil
}
