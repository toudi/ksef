package auth

import (
	"context"
	"errors"
	"ksef/internal/http"
	"ksef/internal/logging"
	baseHttp "net/http"
	"time"
)

const (
	endpointAuthSessions = "/api/v2/auth/sessions"
)

var (
	errNoTokens          = errors.New("sessionTokens == nil")
	errAuthTokenExpired  = errors.New("auth token expired")
	errBothTokensExpired = errors.New("both authToken and refreshToken are expired")
	errUnexpectedStatus  = errors.New("unexpected HTTP response status code")
)

type AuthSessionStatus struct {
	Code        int    `json:"status"`
	Description string `json:"description"`
}
type AuthSession struct {
	ID                     string            `json:"referenceNumber"`
	Current                bool              `json:"isCurrent"`
	StartDate              time.Time         `json:"startDate"`
	AuthMethod             string            `json:"authenticationMethod"`
	Status                 AuthSessionStatus `json:"status"`
	TokenRedeemed          bool              `json:"isTokenRedeemed"`
	RefreshTokenValidUntil time.Time         `json:"refreshTokenValidUntil"`
}
type AuthSessionsResponse struct {
	Sessions []AuthSession `json:"items"`
}

func (tm *TokenManager) getAuthSessionsForToken(
	ctx context.Context, sessionTokens *SessionTokens,
) (sessions *AuthSessionsResponse, err error) {
	if sessionTokens == nil {
		return nil, errNoTokens
	}
	now := time.Now()
	token := sessionTokens.AuthorizationToken
	if token.ExpiresAt.Before(now) {
		return nil, errAuthTokenExpired
	}

	resp, err := tm.httpClient.Request(ctx, http.RequestConfig{
		Method:          baseHttp.MethodGet,
		Headers:         map[string]string{"Authorization": "Bearer " + sessionTokens.AuthorizationToken.Token},
		DestContentType: http.JSON,
		Dest:            &sessions,
	}, endpointAuthSessions)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != baseHttp.StatusOK {
		logging.AuthLogger.Error("nie udało się pobrać listy sesji dla odzyskanego tokenu")
		return nil, errUnexpectedStatus
	}

	return sessions, err
}

func (tm *TokenManager) GetAuthSessions(ctx context.Context) (
	sessions *AuthSessionsResponse, err error,
) {
	if err := tm.restoreTokens(ctx); err != nil {
		return nil, err
	}
	return tm.getAuthSessionsForToken(ctx, tm.sessionTokens)
}
