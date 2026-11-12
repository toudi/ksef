package auth

import (
	"context"
	"errors"
	"sync"
	"time"

	"ksef/core/certsdb"
	"ksef/core/keyring"
	"ksef/internal/client/v2/auth/validator"
	"ksef/internal/http"
	"ksef/internal/logging"
)

const (
	endpointAuthTokenRefresh = "/v2/auth/token/refresh"
)

var ErrTimeoutReadingToken = errors.New("timeout reading authorization token")

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

// Config carries the per-invocation inputs the token manager used to read from
// the CLI configuration layer.
type Config struct {
	NIP string
	// EnvironmentID is the keyring bucket, i.e. the KSeF gateway id. This is used to later
	// construct the keyring key for storing the tokens
	EnvironmentID string
	// HTTPClient is the unauthenticated client used for the auth endpoints.
	HTTPClient *http.Client
	// SkipRestore disables restoring persisted tokens, forcing a fresh
	// authentication.
	SkipRestore bool
}

type TokenManager struct {
	challengeValidator  validator.AuthChallengeValidator
	sessionTokens       *SessionTokens
	mutex               sync.Mutex
	httpClient          *http.Client
	validationReference *validator.ValidationReference
	done                chan struct{}
	ready               chan struct{}
	readyOnce           sync.Once
	obtainNewChallenge  bool
	keyring             keyring.Keyring
	certsDB             *certsdb.CertsDB
	nip                 string
	environmentID       string
	skipRestore         bool
}

func NewTokenManager(
	ctx context.Context,
	cfg Config,
	challengeValidator validator.AuthChallengeValidator,
	keyring keyring.Keyring,
	certsDB *certsdb.CertsDB,
) (*TokenManager, error) {
	if challengeValidator != nil {
		if err := challengeValidator.Initialize(ctx, cfg.HTTPClient); err != nil {
			return nil, err
		}
	}

	tm := &TokenManager{
		httpClient:         cfg.HTTPClient,
		challengeValidator: challengeValidator,
		done:               make(chan struct{}),
		ready:              make(chan struct{}),
		keyring:            keyring,
		certsDB:            certsDB,
		nip:                cfg.NIP,
		environmentID:      cfg.EnvironmentID,
		skipRestore:        cfg.SkipRestore,
	}

	return tm, nil
}

type AuthenticationStatus struct {
	Status struct {
		Code int `json:"code"`
	} `json:"status"`
}

// this is a blocking function that will either:
// - retrieve the most recent authorization token or
// - timeout after 15 seconds
func (t *TokenManager) GetAuthorizationToken(timeout ...time.Duration) (string, error) {
	// a bit of an ugly way to pass a custom tiemout but hey.. it works
	// desperate times call for desperate measures and all that.
	if len(timeout) == 0 {
		timeout = []time.Duration{15 * time.Second}
	}

	deadline := time.Now().Add(timeout[0])
	for time.Now().Before(deadline) {
		t.mutex.Lock()
		sessionTokens := t.sessionTokens
		t.mutex.Unlock()

		if sessionTokens != nil && sessionTokens.AuthorizationToken != nil {
			return sessionTokens.AuthorizationToken.Token, nil
		}

		logging.AuthLogger.Debug("GetAuthorizationToken - awaiting token")

		time.Sleep(500 * time.Millisecond)
	}

	return "", ErrTimeoutReadingToken
}

func (t *TokenManager) currentSessionTokens() *SessionTokens {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	return t.sessionTokens
}

func (t *TokenManager) clearAuthorizationToken() {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.sessionTokens = nil
}

func (t *TokenManager) updateAuthorizationToken(authToken *TokenInfo) {
	logging.AuthLogger.Debug("updateAuthorizationToken()")
	defer logging.AuthLogger.Debug("updateAuthorizationToken() - finish")

	t.mutex.Lock()
	canSafelyPersistTokens := t.sessionTokens != nil

	if t.sessionTokens == nil {
		t.sessionTokens = &SessionTokens{}
	}

	t.sessionTokens.AuthorizationToken = authToken
	t.mutex.Unlock()

	if canSafelyPersistTokens {
		logging.AuthLogger.Debug("updateAuthorizationToken() - try to persist tokens")
		if err := t.persistTokens(); err != nil {
			logging.AuthLogger.Error("unable to persist tokens", "err", err)
		}
		logging.AuthLogger.Debug("updateAuthorizationToken() - tokens successfully persisted")
	}
}

func (t *TokenManager) updateSessionTokens(tokens *SessionTokens, persist bool) {
	logging.AuthLogger.Debug("updateSessionTokens()")
	defer logging.AuthLogger.Debug("updateSessionTokens() - finish")

	t.mutex.Lock()
	t.sessionTokens = tokens
	t.mutex.Unlock()

	if persist {
		logging.AuthLogger.Debug("updateSessionTokens() - try to persist tokens")
		if err := t.persistTokens(); err != nil {
			logging.AuthLogger.Error("unable to persist tokens", "err", err)
		}
		logging.AuthLogger.Debug("updateSessionTokens() - tokens successfully persisted")
	}
}
