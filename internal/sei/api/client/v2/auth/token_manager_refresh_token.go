package auth

import (
	"context"
	"errors"
	"ksef/internal/http"
	baseHttp "net/http"
)

func (t *TokenManager) refreshAccessToken(ctx context.Context, refreshToken string) (*TokenInfo, error) {
	var tokens SessionTokens
	resp, err := t.httpClient.Request(ctx, http.RequestConfig{
		Method:          baseHttp.MethodPost,
		Headers:         map[string]string{"Authorization": "Bearer " + t.sessionTokens.RefreshToken.Token},
		Dest:            &tokens,
		DestContentType: http.JSON,
	}, endpointAuthTokenRefresh)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == baseHttp.StatusOK {
		return tokens.AuthorizationToken, nil
	} else {
		return nil, errors.New("unexpected code")
	}
}
