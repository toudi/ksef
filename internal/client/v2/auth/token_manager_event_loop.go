package auth

import (
	"context"
	"log/slog"
	"time"

	"ksef/internal/client/v2/auth/validator"
	"ksef/internal/logging"
)

// Run drives the auth lifecycle until ctx is cancelled. It is meant to be
// started in its own goroutine by whoever owns the manager.
//
// The manager never closes the keyring: the keyring is constructed and closed
// by the caller (the CLI's resources layer).
func (t *TokenManager) Run(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	defer close(t.done)

	logger := logging.AuthLogger.With("auth", "token manager")

	if !t.skipRestore {
		if err := t.restoreTokens(ctx); err != nil {
			if err != errCannotUseToken {
				logger.Error("fatal (?) error during token restoration", "err", err)
				return
			}
			t.obtainNewChallenge = true
		}
	}

	if t.obtainNewChallenge {
		t.obtainNewChallenge = false
		go t.beginAuth(ctx, logger)
	}

	var lastDebugEvent time.Time

	for {
		select {
		case <-ctx.Done():
			logger.Debug("token manager stopped")
			return

		case now := <-ticker.C:
			// do not be so noisy about loop tick reporting
			if now.Sub(lastDebugEvent) > time.Duration(1*time.Minute) {
				logger.Debug("event loop tick")
				lastDebugEvent = now
			}
			// let's check if we have to refresh the session token
			sessionTokens := t.currentSessionTokens()
			if sessionTokens == nil {
				if t.validationReference == nil {
					logger.Debug("authorization in process")
					break
				}

				if err := t.checkAuthStatus(ctx); err != nil {
					logger.Error("unable to check validation status", "err", err)
				}

				break
			}
			if now.After(sessionTokens.AuthorizationToken.ExpiresAt) {
				logger.Debug("need to refresh authorization token")
				// but only if we're not past refresh token expiration time
				if now.Before(sessionTokens.RefreshToken.ExpiresAt) {
					newToken, err := t.refreshAccessToken(ctx, sessionTokens.RefreshToken.Token)
					if err == nil {
						t.updateAuthorizationToken(newToken)
					}
				} else {
					logger.Warn("too late to refresh token. request new challenge")
					t.clearAuthorizationToken()
					go t.beginAuth(ctx, logger)
				}
			}

		case validatorEvent := <-t.challengeValidator.Event():
			if validatorEvent.State == validator.StateInitialized {
				t.beginAuth(ctx, logger)
			}
			if validatorEvent.State == validator.StateValidationReferenceResult {
				t.validationReference = validatorEvent.ValidationReference
			}
			if validatorEvent.State == validator.StateExit {
				// at the moment it's used only by the "fake" validator of type xadesInit
				// which has a single purpose of dumping the challenge to a file.
				return
			}
		}
	}
}

// Done is closed when Run returns. Callers that started Run can use it to wait
// for a clean shutdown.
func (t *TokenManager) Done() <-chan struct{} {
	return t.done
}

// WaitReady blocks until tokens have been restored or redeemed, or until ctx is
// cancelled. It replaces the old "exit-after-persisting-token" behaviour: the
// caller decides when it has seen enough and cancels its own context.
func (t *TokenManager) WaitReady(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.ready:
		return nil
	}
}

func (t *TokenManager) markReady() {
	t.readyOnce.Do(func() {
		close(t.ready)
	})
}

func (t *TokenManager) beginAuth(ctx context.Context, logger *slog.Logger) {
	logger.Debug("auth handler initialized - receive challenge")
	// validator is ready to validate challenge so we have to obtain the challenge first.
	authChallenge, err := t.getAuthChallenge(ctx)
	if err != nil {
		logger.Error("unable to obtain auth challenge", "err", err)
		return
	}
	if err := t.challengeValidator.ValidateChallenge(ctx, *authChallenge); err != nil {
		logger.Error("unable to validate challenge", "err", err)
	}
}
