package auth

import (
	"context"
	"ksef/internal/http"
	"ksef/internal/logging"
	baseHttp "net/http"
)

const (
	endpointAuthTokenRedeem = "/api/v2/auth/token/redeem"
)

func (t *TokenManager) redeemTokens(ctx context.Context) error {
	// let's fetch the tokens for the initial time
	var tokens SessionTokens

	_, err := t.httpClient.Request(ctx, http.RequestConfig{
		Method:          baseHttp.MethodGet,
		Headers:         map[string]string{"Authorization": "Bearer " + t.validationReference.AuthenticationToken.Token},
		Dest:            &tokens,
		DestContentType: http.JSON,
	}, endpointAuthTokenRedeem)
	if err != nil {
		logging.AuthLogger.Error("unable to redeem token: %w", err)
		return err
	}

	t.updateAuthorizationToken(
		tokens.AuthorizationToken.Token,
		func() {
			t.sessionTokens = &tokens
		},
	)

	return nil
}
