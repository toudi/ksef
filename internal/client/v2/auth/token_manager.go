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
	endpointAuthTokenRefresh     = "/api/v2/auth/token/refresh"
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
	updateChannel       chan TokenUpdate
	mutex               sync.RWMutex
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
		updateChannel:      make(chan TokenUpdate),
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
	ctx, cancel := context.WithTimeout(context.Background(), timeout[0])
	defer cancel()

	// so this temporary channel is for retrieving the *current* token
	tokenChan := make(chan string)
	defer close(tokenChan)

	go t.readToken(tokenChan)

	select {
	case <-ctx.Done():
		return "", ErrTimeoutReadingToken
	case tokenUpdate := <-t.updateChannel:
		return tokenUpdate.Token, tokenUpdate.Err
	case token := <-tokenChan:
		return token, nil
	}
}

// read token from memory. if it exists, push it to tokenChan
func (t *TokenManager) readToken(tokenChan chan string) {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	token := t.sessionTokens
	if token != nil && token.AuthorizationToken != nil {
		tokenChan <- token.AuthorizationToken.Token
	}
}

func (t *TokenManager) updateAuthorizationToken(authToken string, commit func()) {
	logging.AuthLogger.Debug("updateAuthorizationToken()")
	// first send an update to update channel as it doesn't require acquiring a lock
	// which means that the above GetAuthorizationToken function cal capture it
	// in the select loop
	t.updateChannel <- TokenUpdate{
		Token: authToken,
	}

	// t.validationReference = nil

	t.mutex.Lock()
	defer t.mutex.Unlock()

	// callback that has a guarantee of being executed withing mutex lock
	logging.AuthLogger.Debug("updateAuthorizationToken() - call commit function")
	commit()
	logging.AuthLogger.Debug("updateAuthorizationToken() - finish")
}
