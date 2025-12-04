package auth

import (
	"context"
	"ksef/internal/logging"
	"time"
)

func (t *TokenManager) validateSessionTokens(ctx context.Context, sessionTokens *SessionTokens) (canBeUsed bool, refreshed bool) {
	if sessionTokens == nil || sessionTokens.AuthorizationToken == nil || sessionTokens.RefreshToken == nil {
		logging.AuthLogger.Debug("session tokens seem to be empty")
		return false, false
	}

	var now = time.Now()

	if sessionTokens.AuthorizationToken.ExpiresAt.Before(now) && sessionTokens.RefreshToken.ExpiresAt.Before(now) {
		logging.AuthLogger.Debug("both auth token and refresh token are expired")
		return false, false
	}

	// we seem to be ok. let's check if we need to refresh token
	if sessionTokens.AuthorizationToken.ExpiresAt.Before(now) {
		logging.AuthLogger.Debug("auth token is expired. try go obtain a new one")
		// yes, this is the case.
		newToken, err := t.refreshAccessToken(ctx, sessionTokens.RefreshToken.Token)
		if err != nil && err == errUnauthorized {
			logging.AuthLogger.Debug("unable to refresh token")
			return false, false
		}
		sessionTokens.AuthorizationToken = newToken
		return true, true
	}
	logging.AuthLogger.Debug("auth token is not expired. try to validate if it works")
	// we do not have to refresh it so let's try to do a simple GET operation that will resemble a PING
	_, err := t.getAuthSessionsForToken(ctx, sessionTokens)
	return err == nil, false
}
