package auth

import (
	"context"
	"ksef/internal/http"
	"ksef/internal/logging"
	baseHttp "net/http"
)

const (
	endpointAuthTokenRedeem = "/v2/auth/token/redeem"
)

func (t *TokenManager) redeemTokens(ctx context.Context) error {
	// let's fetch the tokens for the initial time
	var tokens SessionTokens

	_, err := t.httpClient.Request(ctx, http.RequestConfig{
		Method:          baseHttp.MethodPost,
		Headers:         map[string]string{"Authorization": "Bearer " + t.validationReference.AuthenticationToken.Token},
		Dest:            &tokens,
		DestContentType: http.JSON,
		ExpectedStatus:  baseHttp.StatusOK,
	}, endpointAuthTokenRedeem)
	if err != nil {
		logging.AuthLogger.Error("unable to redeem token", "err", err)
		return err
	}

	t.updateSessionTokens(&tokens, true)

	if t.vip.GetBool(FlagExitAfterPersistingToken) {
		t.finished = true
	}

	return nil
}
