package kseftoken

import (
	"ksef/internal/config"
	"ksef/internal/logging"
	"ksef/internal/sei/api/client/v2/auth/validator"

	"github.com/zalando/go-keyring"
)

type KsefTokenHandler struct {
	nip          string
	ksefToken    string // just to distinguish it from the session token
	apiConfig    config.APIConfig
	eventChannel chan validator.AuthEvent
	finished     bool
}

func NewKsefTokenHandler(apiConfig config.APIConfig, nip string) validator.AuthChallengeValidator {
	var err error

	validator := &KsefTokenHandler{
		nip:          nip,
		apiConfig:    apiConfig,
		eventChannel: make(chan validator.AuthEvent),
	}

	// let's try to retrieve it from keyring
	ksefToken, err := retrieveKsefTokenFromKeyring(apiConfig.Host, nip)
	if err != nil {
		// that's not a fatal error because the program also supports overriding the token directly
		logging.AuthLogger.Warn("unable to retrieve KSeF token from keyring")
	} else {
		validator.SetKsefToken(ksefToken)
	}

	return validator
}

func (kt *KsefTokenHandler) Event() chan validator.AuthEvent {
	return kt.eventChannel
}

func (kt *KsefTokenHandler) SetKsefToken(token string) {
	kt.ksefToken = token
	kt.eventChannel <- validator.AuthEvent{
		State: validator.StateInitialized,
	}

}

func retrieveKsefTokenFromKeyring(gateway string, issuerNip string) (string, error) {
	return keyring.Get(gateway, issuerNip)
}

func PersistKsefTokenToKeyring(gateway string, issuerNip string, token string) error {
	return keyring.Set(gateway, issuerNip, token)
}
