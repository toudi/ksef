package auth

import (
	"context"
	"ksef/internal/logging"
	"ksef/internal/sei/api/client/v2/auth/validator"
	"time"
)

func (t *TokenManager) Stop() {
	t.finished = true
}

func (t *TokenManager) Done() chan struct{} {
	return t.done
}

func (t *TokenManager) Run() {
	ticker := time.NewTicker(1 * time.Second)

	ctx := context.Background()
	logger := logging.AuthLogger.With("auth", "token manager")

	for !t.finished {
		select {
		case now := <-ticker.C:
			logger.Debug("event loop tick")
			// let's check if we have to refresh the session token
			if t.sessionTokens == nil {
				if t.validationReference == nil {
					logger.Debug("authorization in process")
					break
				}

				if err := t.checkAuthStatus(ctx); err != nil {
					logger.Error("unable to check validation status", "err", err)
				}

				break
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
					logger.Error("too late to refresh token")
				}
			}

		case validatorEvent := <-t.challengeValidator.Event():
			if validatorEvent.State == validator.StateInitialized {
				logger.Debug("auth handler initialized - receive challenge")
				// validator is ready to validate challenge so we have to obtain the challenge first.
				authChallenge, err := t.getAuthChallenge(ctx)
				if err != nil {
					logger.Error("unable to obtain auth challenge", "err", err)
					break
				}
				if err := t.challengeValidator.ValidateChallenge(ctx, *authChallenge); err != nil {
					logger.Error("unable to validate challenge", "err", err)
				}
			}
			if validatorEvent.State == validator.StateValidationReferenceResult {
				t.validationReference = validatorEvent.ValidationReference
			}
			if validatorEvent.State == validator.StateExit {
				// at the moment it's used only by the "fake" validator of type xadesInit
				// which has a single purpose of dumping the challenge to a file.
				t.finished = true
			}
		}
	}

	ticker.Stop()

	t.done <- struct{}{}
}
