package xades

import (
	"bytes"
	"context"
	"ksef/internal/certsdb"
	"ksef/internal/config"
	"ksef/internal/http"
	"ksef/internal/sei/api/client/v2/auth/validator"
)

type AuthHandler struct {
	nip          string
	certificate  certsdb.Certificate
	eventChannel chan validator.AuthEvent
	httpClient   *http.Client
}

func NewAuthHandler(apiConfig config.APIConfig, nip string, certificate certsdb.Certificate) validator.AuthChallengeValidator {
	handler := &AuthHandler{
		eventChannel: make(chan validator.AuthEvent),
		nip:          nip,
		certificate:  certificate,
	}

	return handler
}

func (e *AuthHandler) Event() chan validator.AuthEvent {
	return e.eventChannel
}

func (e *AuthHandler) initialize() {
	e.eventChannel <- validator.AuthEvent{
		State: validator.StateInitialized,
	}
}

func (e *AuthHandler) Initialize(httpClient *http.Client) {
	e.httpClient = httpClient

	go e.initialize()
}

func (e *AuthHandler) ValidateChallenge(ctx context.Context, challenge validator.AuthChallenge) error {
	// we have our challenge. we now need to sign it and send using validateSignedChallenge
	var sourceDocument bytes.Buffer
	err := dumpChallengeToWriter(challenge, e.nip, &sourceDocument)
	if err != nil {
		return err
	}
	// great. now we can sign it using the certificate
	var signedDocument bytes.Buffer
	if err = SignAuthChallenge(&sourceDocument, e.certificate, &signedDocument); err != nil {
		return err
	}
	// perfect. final step - let's post it to the validation endpoint
	return validateSignedChallenge(
		ctx,
		e.httpClient,
		bytes.NewReader(signedDocument.Bytes()),
		func(resp validator.ValidationReference) {
			e.eventChannel <- validator.AuthEvent{
				State:               validator.StateValidationReferenceResult,
				ValidationReference: &resp,
			}
		},
	)
}
