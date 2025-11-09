package auth

import (
	"context"
	"ksef/internal/config"
	"ksef/internal/http"
	"ksef/internal/logging"
	baseHttp "net/http"

	"github.com/spf13/viper"
)

const (
	endpointAuthTokenRedeem = "/api/v2/auth/token/redeem"
)

func (t *TokenManager) redeemTokens(ctx context.Context) error {
	// let's fetch the tokens for the initial time
	var tokens SessionTokens

	_, err := t.httpClient.Request(ctx, http.RequestConfig{
		Method:          baseHttp.MethodPost,
		Headers:         map[string]string{"Authorization": "Bearer " + t.validationReference.AuthenticationToken.Token},
		Dest:            &tokens,
		DestContentType: http.JSON,
	}, endpointAuthTokenRedeem)
	if err != nil {
		logging.AuthLogger.Error("unable to redeem token", "err", err)
		return err
	}

	t.updateAuthorizationToken(
		tokens.AuthorizationToken.Token,
		func() {
			vip := viper.GetViper()

			t.sessionTokens = &tokens
			nip, err := config.GetNIP(vip)
			if err != nil {
				logging.AuthLogger.Error("unable to validate NIP", "err", err)
			}
			if err := t.PersistTokens(config.GetGateway(vip), nip); err != nil {
				logging.AuthLogger.Error("unable to persist tokens", "err", err)
			}
		},
	)

	return nil
}
