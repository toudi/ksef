package token

import (
	_ "embed"
	"io"
	"ksef/internal/client/v2/auth/validator"
	"text/template"
)

var (
	//go:embed "templates/AuthTokenRequest.xml"
	challengeRequest         string
	challengeRequestTemplate *template.Template
)

type authChallengeRequestVars struct {
	Challenge  string
	SubjectNIP string
}

func init() {
	challengeRequestTemplate, _ = template.New("challenge").Parse(challengeRequest)
}

func dumpChallengeToWriter(challenge validator.AuthChallenge, nip string, dest io.Writer) error {
	return challengeRequestTemplate.Execute(
		dest,
		authChallengeRequestVars{
			Challenge:  challenge.Challenge,
			SubjectNIP: nip,
		},
	)
}
