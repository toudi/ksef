package auth

import (
	"context"
	httpClient "ksef/internal/http"
	"net/http"
	"time"
)

const endpointAuthChallenge = "/api/v2/auth/challenge"

type authChallengeResponse struct {
	Challenge string    `json:"challenge"`
	Timestamp time.Time `json:"timestamp"`
}

func (ah AuthHandler) ObtainSessionToken(ctx context.Context) (*ValidationResponse, error) {
	var authResponse authChallengeResponse

	_, err := ah.httpClient.Request(
		ctx, httpClient.RequestConfig{
			Method:          http.MethodPost,
			ContentType:     httpClient.JSON,
			Dest:            &authResponse,
			DestContentType: httpClient.JSON,
			ExpectedStatus:  http.StatusOK,
		},
		endpointAuthChallenge,
	)

	if err != nil {
		return nil, err
	}

	return ah.challengeValidator.ValidateChallenge(ctx, authResponse)
}
