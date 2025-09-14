package auth

import (
	"context"
	"fmt"
	"ksef/internal/logging"
	"ksef/internal/sei/api/client/v2/auth/validator"
	"time"
)

func (t *TokenManager) Stop() {
	t.finished = true
}

func (t *TokenManager) Run() {
	ticker := time.NewTicker(1 * time.Second)

	ctx := context.Background()

	for !t.finished {
		select {
		case now := <-ticker.C:
			// let's check if we have to refresh the session token
			if t.sessionTokens == nil {
				if t.validationReference == nil {
					break
				}

				if err := t.checkAuthStatus(ctx); err != nil {
					logging.AuthLogger.Error("unable to check validation status: %w", err)
				}
			}
			if now.After(t.sessionTokens.AuthorizationToken.ExpiresAt) {
				// but only if we're not past refresh token expiration time
				if now.Before(t.sessionTokens.RefreshToken.ExpiresAt) {
					newToken, err := t.refreshAccessToken(ctx, t.sessionTokens.RefreshToken.Token)
					if err == nil {
						t.updateAuthorizationToken(newToken.Token, func() {
							t.sessionTokens.AuthorizationToken = newToken
						})
					}
				} else {
					fmt.Printf("too late to refresh token")
				}
			}

		case validatorEvent := <-t.challengeValidator.Event():
			if validatorEvent.State == validator.StateInitialized {
				// validator is ready to validate challenge so we have to obtain the challenge first.
				authChallenge, err := t.getAuthChallenge(ctx)
				if err != nil {
					logging.AuthLogger.Error("unable to obtain auth challenge: %w", err)
					break
				}
				validationRef, err := t.challengeValidator.ValidateChallenge(ctx, t.httpClient, *authChallenge)
				if err != nil {
					logging.AuthLogger.Error("unable to validate challenge: %w", err)
				} else {
					t.validationReference = validationRef
				}
				break
			}
		}
	}

	ticker.Stop()
}
