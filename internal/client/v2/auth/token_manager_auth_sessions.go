package auth

import (
	"context"
	"errors"
	baseHttp "net/http"
	"time"

	"ksef/core/ksef"
	"ksef/internal/http"
	"ksef/internal/logging"
)

const (
	endpointAuthSessions = "/v2/auth/sessions"
)

var (
	errNoTokens          = errors.New("sessionTokens == nil")
	errAuthTokenExpired  = errors.New("auth token expired")
	errBothTokensExpired = errors.New("both authToken and refreshToken are expired")
	errUnexpectedStatus  = errors.New("unexpected HTTP response status code")
)

func (tm *TokenManager) getAuthSessionsForToken(
	ctx context.Context, sessionTokens *SessionTokens,
) (sessions *ksef.AuthSessions, err error) {
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
	sessions *ksef.AuthSessions, err error,
) {
	if err := tm.restoreTokens(ctx); err != nil {
		return nil, err
	}
	return tm.getAuthSessionsForToken(ctx, tm.currentSessionTokens())
}
