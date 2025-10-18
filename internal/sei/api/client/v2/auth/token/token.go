package token

import (
	"context"
	"ksef/internal/config"
	"ksef/internal/environment"
	"ksef/internal/http"
	"ksef/internal/sei/api/client/v2/auth/validator"

	"github.com/zalando/go-keyring"
)

type TokenHandler struct {
	env          environment.Environment
	nip          string
	eventChannel chan validator.AuthEvent
}

func NewAuthHandler(apiConfig config.APIConfig, nip string) validator.AuthChallengeValidator {
	handler := &TokenHandler{
		env:          apiConfig.Environment.Environment,
		eventChannel: make(chan validator.AuthEvent),
		nip:          nip,
	}

	return handler
}

func (e *TokenHandler) Event() chan validator.AuthEvent {
	return e.eventChannel
}

func (e *TokenHandler) initialize() {
	sessionTokens, err := keyring.Get(string(e.env)+"-sessionTokens", e.nip)
	if err != nil {
		panic("unable to retrieve token")
	}
	e.eventChannel <- validator.AuthEvent{
		State:         validator.StateTokensReady,
		SessionTokens: sessionTokens,
	}
}

func (e *TokenHandler) Initialize(httpClient *http.Client) {
	go e.initialize()
}

func (e *TokenHandler) ValidateChallenge(ctx context.Context, challenge validator.AuthChallenge) error {
	return nil
}
