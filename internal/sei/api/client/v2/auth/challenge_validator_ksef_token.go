package auth

import (
	"context"
	"fmt"
	"ksef/internal/encryption"
	httpClient "ksef/internal/http"
	"time"
)

const endpointValidateKsefToken = "/api/v2/auth/ksef-token"
const contextIdentifierTypeNIP = "Nip"

type ksefTokenAuthValidator struct {
	nip     string
	token   string
	handler *AuthHandler
}

type contextIdentifier struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type validationRequest struct {
	Challenge      string            `json:"challenge"`
	Identifier     contextIdentifier `json:"contextIdentifier"`
	EncryptedToken string            `json:"encryptedToken"`
}

func NewKsefTokenAuthValidator(token string) AuthValidator {
	return &ksefTokenAuthValidator{}
}

func (kt *ksefTokenAuthValidator) encryptToken(tokenPlaintext string, timestamp time.Time) (string, error) {
	encryptedBytes, err := encryption.EncryptMessageWithCertificate(
		kt.handler.config.Certificate.PEM(),
		fmt.Appendf([]byte{}, "%s|%d", tokenPlaintext, timestamp.UnixMilli()),
	)
	if err != nil {
		return "", err
	}
	return string(encryptedBytes), nil
}

func (kt *ksefTokenAuthValidator) ValidateChallenge(ctx context.Context, challenge authChallengeResponse) (*ValidationResponse, error) {
	var err error

	var body = validationRequest{
		Challenge: challenge.Challenge,
		Identifier: contextIdentifier{
			Type:  contextIdentifierTypeNIP,
			Value: kt.nip,
		},
	}

	if body.EncryptedToken, err = kt.encryptToken(kt.token, challenge.Timestamp); err != nil {
		return nil, err
	}

	var resp ValidationResponse

	_, err = kt.handler.httpClient.Request(
		ctx,
		httpClient.RequestConfig{
			Body:            body,
			ContentType:     httpClient.JSON,
			Dest:            &resp,
			DestContentType: httpClient.JSON,
		},
		endpointValidateKsefToken,
	)

	return &resp, err
}
