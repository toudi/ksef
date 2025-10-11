package xades

import (
	"context"
	"ksef/internal/config"
	"ksef/internal/http"
	"ksef/internal/sei/api/client/v2/auth/validator"
	"os"
)

// this auth handler is a special one - it's only purpose is to dump
// the content of auth challenge to a file. This way one can either
// sign it with the trusted profile, or manually sign it with xades sign command
type AuthChallengeDumper struct {
	nip          string
	eventChannel chan validator.AuthEvent
	httpClient   *http.Client
	outputFile   string
}

func NewChallengeDumperHandler(apiConfig config.APIConfig, nip string, outputFile string) validator.AuthChallengeValidator {
	handler := &AuthChallengeDumper{
		eventChannel: make(chan validator.AuthEvent),
		nip:          nip,
		outputFile:   outputFile,
	}

	return handler
}

func (e *AuthChallengeDumper) initialize() {
	e.eventChannel <- validator.AuthEvent{
		State: validator.StateInitialized,
	}
}

func (e *AuthChallengeDumper) Initialize(httpClient *http.Client) {
	e.httpClient = httpClient

	go e.initialize()
}

func (e *AuthChallengeDumper) Event() chan validator.AuthEvent {
	return e.eventChannel
}

func (e *AuthChallengeDumper) ValidateChallenge(ctx context.Context, challenge validator.AuthChallenge) error {
	return e.dumpChallenge(challenge)
}

type authChallengeRequestVars struct {
	Challenge  string
	SubjectNIP string
}

func (e *AuthChallengeDumper) dumpChallenge(challenge validator.AuthChallenge) error {
	challengeFile, err := os.Create("AuthTokenRequest.xml")
	if err != nil {
		return err
	}

	defer challengeFile.Close()

	err = challengeRequestTemplate.Execute(
		challengeFile,
		authChallengeRequestVars{
			Challenge:  challenge.Challenge,
			SubjectNIP: e.nip,
		},
	)
	if err != nil {
		return err
	}

	go func() {
		e.eventChannel <- validator.AuthEvent{
			State: validator.StateExit,
		}
	}()

	return nil
}
