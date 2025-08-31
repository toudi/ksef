package auth

import (
	"context"
	"errors"
	"fmt"
	"ksef/internal/http"
	baseHttp "net/http"
	"time"
)

const (
	endpointAuthStatus       = "/api/v2/auth/%s"
	endpointAuthTokenRedeem  = "/api/v2/auth/token/redeem"
	endpointAuthTokenRefresh = "/api/v2/auth/token/refresh"
	authStatusCodeSuccess    = 200
)

type TokenUpdate struct {
	Token string
	Err   error
}

type TokenInfo struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"validUntil"`
}
type SessionTokens struct {
	AuthorizationToken *TokenInfo `json:"accessToken"`
	RefreshToken       *TokenInfo `json:"refreshToken"`
}
type TokenManager struct {
	finished            bool
	authenticationToken string
	sessionTokens       *SessionTokens
	C                   chan TokenUpdate
}

type AuthenticationStatus struct {
	Status struct {
		Code int `json:"code"`
	} `json:"status"`
}

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
					t.sessionTokens = &tokens
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
						t.sessionTokens.AuthorizationToken = newToken
					}
				} else {
					fmt.Printf("too late to refresh token")
				}
			}
		}
	}

	ticker.Stop()
}

func (t *TokenManager) GetAccessToken(ctx context.Context, httpClient http.Client, refreshToken string) (*TokenInfo, error) {
	var tokens SessionTokens
	resp, err := httpClient.Request(ctx, http.RequestConfig{
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
