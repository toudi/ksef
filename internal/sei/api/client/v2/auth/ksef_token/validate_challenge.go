package kseftoken

import (
	"context"
	"fmt"
	"ksef/internal/certsdb"
	"ksef/internal/encryption"
	"ksef/internal/http"
	"ksef/internal/sei/api/client/v2/auth/validator"
	"time"
)

const (
	endpointValidateKsefToken = "/api/v2/auth/ksef-token"
	contextIdentifierTypeNIP  = "Nip"
)

type contextIdentifier struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type validationRequest struct {
	Challenge      string            `json:"challenge"`
	Identifier     contextIdentifier `json:"contextIdentifier"`
	EncryptedToken string            `json:"encryptedToken"`
}

func (kt *KsefTokenHandler) ValidateChallenge(ctx context.Context, challenge validator.AuthChallenge) error {
	var err error

	var body = validationRequest{
		Challenge: challenge.Challenge,
		Identifier: contextIdentifier{
			Type:  contextIdentifierTypeNIP,
			Value: kt.nip,
		},
	}

	if body.EncryptedToken, err = kt.encryptToken(kt.ksefToken, challenge.Timestamp); err != nil {
		return err
	}

	var resp validator.ValidationReference

	_, err = kt.httpClient.Request(
		ctx,
		http.RequestConfig{
			Body:            body,
			ContentType:     http.JSON,
			Dest:            &resp,
			DestContentType: http.JSON,
		},
		endpointValidateKsefToken,
	)

	if err != nil {
		return err
	}

	go func() {
		kt.eventChannel <- validator.AuthEvent{
			State:               validator.StateValidationReferenceResult,
			ValidationReference: &resp,
		}
	}()

	return nil
}

func (kt *KsefTokenHandler) encryptToken(tokenPlaintext string, timestamp time.Time) (string, error) {
	certificate, err := kt.apiConfig.CertificatesDB.GetByUsage(certsdb.UsageTokenEncryption, "")
	if err != nil {
		return "", err
	}
	encryptedBytes, err := encryption.EncryptMessageWithCertificate(
		certificate.Filename(),
		fmt.Appendf([]byte{}, "%s|%d", tokenPlaintext, timestamp.UnixMilli()),
	)
	if err != nil {
		return "", err
	}
	return string(encryptedBytes), nil
}
