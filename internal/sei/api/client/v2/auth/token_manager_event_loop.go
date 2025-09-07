package auth

import (
	"context"
	"fmt"
	"ksef/internal/http"
	baseHttp "net/http"
	"time"
)

const (
	endpointAuthTokenRedeem = "/api/v2/auth/token/redeem"
)

func (t *TokenManager) Stop() {
	t.finished = true
}

func (t *TokenManager) Run() {
	ticker := time.NewTicker(1 * time.Second)

	var httpClient http.Client
	ctx := context.Background()

	for !t.finished {
		<-ticker.C
		if t.sessionTokens == nil {
			if t.authenticationToken != "" {
				var authHeaders = map[string]string{"Authorization": "Bearer " + t.authenticationToken}
				// session tokens are empty - we have to check if they can be retrieved
				var authStatus AuthenticationStatus
				_, _ = httpClient.Request(ctx, http.RequestConfig{
					Headers:         authHeaders,
					Dest:            &authStatus,
					DestContentType: http.JSON,
					ExpectedStatus:  baseHttp.StatusOK,
					Method:          baseHttp.MethodPost,
				}, fmt.Sprintf(endpointAuthStatus, t.authenticationToken))
				if authStatus.Status.Code == authStatusCodeSuccess {
					// let's fetch the tokens for the initial time
					var tokens SessionTokens
					_, _ = httpClient.Request(ctx, http.RequestConfig{
						Method:          baseHttp.MethodGet,
						Headers:         authHeaders,
						Dest:            &tokens,
						DestContentType: http.JSON,
					}, endpointAuthTokenRedeem)
					t.updateAuthorizationToken(
						tokens.AuthorizationToken.Token,
						func() {
							t.sessionTokens = &tokens
						},
					)
				}
			}
		} else {
			// if we do have session tokens, let's check if we need to refresh them
			var now = time.Now()
			if now.After(t.sessionTokens.AuthorizationToken.ExpiresAt) {
				// but only if we're not past refresh token expiration time
				if now.Before(t.sessionTokens.RefreshToken.ExpiresAt) {
					newToken, err := t.GetAccessToken(ctx, httpClient, t.sessionTokens.RefreshToken.Token)
					if err == nil {
						t.updateAuthorizationToken(newToken.Token, func() {
							t.sessionTokens.AuthorizationToken = newToken
						})
					}
				} else {
					fmt.Printf("too late to refresh token")
				}
			}
		}
	}

	ticker.Stop()
}
