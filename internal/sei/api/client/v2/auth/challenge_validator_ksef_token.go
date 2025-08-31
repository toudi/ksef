package auth

import (
	"context"
	"fmt"
	"ksef/internal/config"
	"ksef/internal/encryption"
	HTTP "ksef/internal/http"
	"ksef/internal/logging"
	"time"

	"github.com/zalando/go-keyring"
)

const endpointValidateKsefToken = "/api/v2/auth/ksef-token"
const contextIdentifierTypeNIP = "Nip"

type ksefTokenAuthValidator struct {
	nip       string
	ksefToken string // just to distinguish it from the session token
	apiConfig config.APIConfig
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

func NewKsefTokenAuthValidator(apiConfig config.APIConfig, nip string) AuthValidator {
	var err error

	validator := &ksefTokenAuthValidator{
		nip:       nip,
		apiConfig: apiConfig,
	}

	// let's try to retrieve it from keyring
	if validator.ksefToken, err = retrieveKsefTokenFromKeyring(apiConfig.Host, nip); err != nil {
		// that's not a fatal error because the program also supports overriding the token directly
		logging.AuthLogger.Warn("unable to retrieve KSeF token from keyring")
	}

	return validator
}

func (kt *ksefTokenAuthValidator) SetKsefToken(token string) {
	kt.ksefToken = token
}

func (kt *ksefTokenAuthValidator) encryptToken(tokenPlaintext string, timestamp time.Time) (string, error) {
	encryptedBytes, err := encryption.EncryptMessageWithCertificate(
		kt.apiConfig.Certificate.PEM(),
		fmt.Appendf([]byte{}, "%s|%d", tokenPlaintext, timestamp.UnixMilli()),
	)
	if err != nil {
		return "", err
	}
	return string(encryptedBytes), nil
}

func (kt *ksefTokenAuthValidator) ValidateChallenge(ctx context.Context, httpClient HTTP.Client, challenge authChallengeResponse) (*ValidationResponse, error) {
	var err error

	var body = validationRequest{
		Challenge: challenge.Challenge,
		Identifier: contextIdentifier{
			Type:  contextIdentifierTypeNIP,
			Value: kt.nip,
		},
	}

	if body.EncryptedToken, err = kt.encryptToken(kt.ksefToken, challenge.Timestamp); err != nil {
		return nil, err
	}

	var resp ValidationResponse

	_, err = httpClient.Request(
		ctx,
		HTTP.RequestConfig{
			Body:            body,
			ContentType:     HTTP.JSON,
			Dest:            &resp,
			DestContentType: HTTP.JSON,
		},
		endpointValidateKsefToken,
	)

	return &resp, err
}

func retrieveKsefTokenFromKeyring(gateway string, issuerNip string) (string, error) {
	return keyring.Get(gateway, issuerNip)
}

func PersistKsefTokenToKeyring(gateway string, issuerNip string, token string) error {
	return keyring.Set(gateway, issuerNip, token)
}
