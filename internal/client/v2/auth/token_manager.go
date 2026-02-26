package auth

import (
	"context"
	"errors"
	"ksef/internal/certsdb"
	"ksef/internal/client/v2/auth/validator"
	"ksef/internal/http"
	"ksef/internal/keyring"
	"ksef/internal/logging"
	"ksef/internal/runtime"
	"sync"
	"time"

	"github.com/spf13/viper"
)

const (
	endpointAuthTokenRefresh     = "/v2/auth/token/refresh"
	FlagDoNotRestoreTokens       = "no-restore-tokens"
	FlagExitAfterPersistingToken = "exit-after-persisting-token"
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

type TokenManager struct {
	challengeValidator  validator.AuthChallengeValidator
	finished            bool
	sessionTokens       *SessionTokens
	mutex               sync.Mutex
	httpClient          *http.Client
	validationReference *validator.ValidationReference
	done                chan struct{}
	vip                 *viper.Viper
	obtainNewChallenge  bool
	keyring             keyring.Keyring
	certsDB             *certsdb.CertificatesDB
}

func NewTokenManager(
	ctx context.Context,
	vip *viper.Viper,
	challengeValidator validator.AuthChallengeValidator,
) (*TokenManager, error) {
	environment := runtime.GetEnvironment(vip)
	httpClient := http.NewClient(environment.API)

	if challengeValidator != nil {
		if err := challengeValidator.Initialize(ctx, httpClient); err != nil {
			return nil, err
		}
	}

	keyring, err := keyring.NewKeyring(vip)
	if err != nil {
		logging.AuthLogger.Error("unable to initialize keyring", "err", err)
		return nil, err
	}
	certsDB, err := certsdb.OpenOrCreate(vip)
	if err != nil {
		return nil, err
	}

	tm := &TokenManager{
		httpClient:         httpClient,
		challengeValidator: challengeValidator,
		done:               make(chan struct{}),
		vip:                vip,
		keyring:            keyring,
		certsDB:            certsDB,
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

func (t *TokenManager) updateAuthorizationToken(authToken *TokenInfo) {
	logging.AuthLogger.Debug("updateAuthorizationToken()")
	defer logging.AuthLogger.Debug("updateAuthorizationToken() - finish")

	t.mutex.Lock()
	defer t.mutex.Unlock()

	canSafelyPersistTokens := t.sessionTokens != nil

	if t.sessionTokens == nil {
		t.sessionTokens = &SessionTokens{}
	}

	t.sessionTokens.AuthorizationToken = authToken

	if canSafelyPersistTokens {
		logging.AuthLogger.Debug("updateAuthorizationToken() - try to persist tokens")
		if err := t.persistTokens(); err != nil {
			logging.AuthLogger.Error("unable to persist tokens", "err", err)
		}
		logging.AuthLogger.Debug("updateAuthorizationToken() - tokens successfully persisted")
	}
}

func (t *TokenManager) updateSessionTokens(tokens *SessionTokens) {
	logging.AuthLogger.Debug("updateSessionTokens()")
	defer logging.AuthLogger.Debug("updateSessionTokens() - finish")

	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.sessionTokens = tokens

	logging.AuthLogger.Debug("updateSessionTokens() - try to persist tokens")
	if err := t.persistTokens(); err != nil {
		logging.AuthLogger.Error("unable to persist tokens", "err", err)
	}
	logging.AuthLogger.Debug("updateSessionTokens() - tokens successfully persisted")
}
