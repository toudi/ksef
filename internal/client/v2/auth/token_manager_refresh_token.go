package auth

import (
	"context"
	"errors"
	"ksef/internal/http"
	baseHttp "net/http"
)

var errUnauthorized = errors.New("unauthorized - token is probably deactivated")

func (t *TokenManager) refreshAccessToken(ctx context.Context, refreshToken string) (*TokenInfo, error) {
	var tokens SessionTokens
	resp, err := t.httpClient.Request(ctx, http.RequestConfig{
		Method:          baseHttp.MethodPost,
		Headers:         map[string]string{"Authorization": "Bearer " + refreshToken},
		Dest:            &tokens,
		DestContentType: http.JSON,
	}, endpointAuthTokenRefresh)
	if err != nil {
		return nil, err
	}
	switch resp.StatusCode {
	case baseHttp.StatusOK:
		return tokens.AuthorizationToken, nil
	case baseHttp.StatusUnauthorized:
		return nil, errUnauthorized
	default:
		return nil, errors.New("unexpected code")
	}
}
