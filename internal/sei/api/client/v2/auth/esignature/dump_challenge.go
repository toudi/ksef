package esignature

import (
	"ksef/internal/sei/api/client/v2/auth/validator"
	"os"
)

type authChallengeRequestVars struct {
	Challenge  string
	SubjectNIP string
}

func (e *eSignatureTokenHandler) dumpChallenge(challenge validator.AuthChallenge) error {
	challengeFile, err := os.Create("AuthTokenRequest.xml")
	if err != nil {
		return err
	}

	defer challengeFile.Close()

	err = e.authTokenRequestTemplate.Execute(
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
