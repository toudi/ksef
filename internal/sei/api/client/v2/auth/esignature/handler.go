package esignature

import (
	"context"
	_ "embed"
	"ksef/internal/config"
	"ksef/internal/http"
	"ksef/internal/sei/api/client/v2/auth/validator"
	"text/template"
)

//go:embed "AuthTokenRequest.xml"
var authTokenRequestSource string

type eSignatureTokenHandler struct {
	nip                      string
	signedFile               string
	eventChannel             chan validator.AuthEvent
	authTokenRequestTemplate *template.Template
	httpClient               *http.Client
}

func NewESignatureTokenHandler(apiConfig config.APIConfig, nip string, signedFile string) validator.AuthChallengeValidator {
	// load the signature template to memory
	tmpl, err := template.New("AuthTokenRequest.xml").Parse(authTokenRequestSource)
	if err != nil {
		panic("unable to parse embedded AuthTokenRequest.xml")
	}

	handler := &eSignatureTokenHandler{
		eventChannel:             make(chan validator.AuthEvent),
		authTokenRequestTemplate: tmpl,
		nip:                      nip,
		signedFile:               signedFile,
	}

	return handler
}

func (e *eSignatureTokenHandler) initialize() {
	e.eventChannel <- validator.AuthEvent{
		State: validator.StateInitialized,
	}
}

func (e *eSignatureTokenHandler) Initialize(httpClient *http.Client) {
	e.httpClient = httpClient

	if e.nip != "" {
		go e.initialize()
	} else {
		go e.validateSignedChallenge(e.signedFile)
	}
}

func (e *eSignatureTokenHandler) Event() chan validator.AuthEvent {
	return e.eventChannel
}

func (e *eSignatureTokenHandler) ValidateChallenge(ctx context.Context, challenge validator.AuthChallenge) error {
	if e.nip != "" {
		// in this mode we're only going to dump the request to a file
		return e.dumpChallenge(challenge)
	}
	return nil
}
