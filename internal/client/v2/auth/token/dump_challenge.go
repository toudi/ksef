package token

import (
	_ "embed"
	"io"
	"ksef/internal/certsdb"
	"ksef/internal/client/v2/auth/validator"
	"text/template"
)

var (
	//go:embed "templates/AuthTokenRequest-nip.xml"
	challengeRequestSubjectNIP         string
	challengeRequestTemplateSubjectNIP *template.Template
	//go:embed "templates/AuthTokenRequest-internal-id.xml"
	challengeRequestSubjectInternalId         string
	challengeRequestTemplateSubjectInternalId *template.Template
)

type authChallengeRequestVars struct {
	Challenge  string
	SubjectNIP string
	InternalId *string
}

func init() {
	challengeRequestTemplateSubjectNIP, _ = template.New("context-nip").Parse(challengeRequestSubjectNIP)
	challengeRequestTemplateSubjectInternalId, _ = template.New("context-internal-id").Parse(challengeRequestSubjectInternalId)
}

func dumpChallengeToWriter(challenge validator.AuthChallenge, nip string, certificate certsdb.Certificate, dest io.Writer) error {
	internalId := certificate.GetInternalIDForNIP(nip)
	if internalId != nil {
		return challengeRequestTemplateSubjectInternalId.Execute(
			dest,
			authChallengeRequestVars{
				Challenge:  challenge.Challenge,
				InternalId: internalId,
			},
		)
	}
	return challengeRequestTemplateSubjectNIP.Execute(
		dest,
		authChallengeRequestVars{
			Challenge:  challenge.Challenge,
			SubjectNIP: nip,
		},
	)
}
