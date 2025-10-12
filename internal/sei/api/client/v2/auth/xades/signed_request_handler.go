package xades

import (
	"context"
	_ "embed"
	"ksef/internal/config"
	"ksef/internal/http"
	"ksef/internal/sei/api/client/v2/auth/validator"
	"os"
)

type signedRequestHandler struct {
	signedFile   string
	eventChannel chan validator.AuthEvent
	httpClient   *http.Client
}

func NewSignedRequestHandler(apiConfig config.APIConfig, signedFile string) validator.AuthChallengeValidator {
	handler := &signedRequestHandler{
		eventChannel: make(chan validator.AuthEvent),
		signedFile:   signedFile,
	}

	return handler
}

func (e *signedRequestHandler) initialize() {
	e.eventChannel <- validator.AuthEvent{
		State: validator.StateInitialized,
	}
}

func (e *signedRequestHandler) Initialize(httpClient *http.Client) {
	e.httpClient = httpClient

	signedChallenge, err := os.Open(e.signedFile)
	if err != nil {
		panic(err)
	}

	go validateSignedChallenge(
		context.Background(),
		e.httpClient,
		signedChallenge,
		func(resp validator.ValidationReference) {
			e.eventChannel <- validator.AuthEvent{
				State:               validator.StateValidationReferenceResult,
				ValidationReference: &resp,
			}
		},
	)
}

func (e *signedRequestHandler) Event() chan validator.AuthEvent {
	return e.eventChannel
}

func (e *signedRequestHandler) ValidateChallenge(ctx context.Context, challenge validator.AuthChallenge) error {
	return nil
}
