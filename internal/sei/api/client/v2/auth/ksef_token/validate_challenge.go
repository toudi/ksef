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

func (kt *KsefTokenHandler) ValidateChallenge(ctx context.Context, httpClient *http.Client, challenge validator.AuthChallenge) (*validator.ValidationReference, error) {
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

	var resp validator.ValidationReference

	_, err = httpClient.Request(
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
		return nil, err
	}

	return &resp, nil
}

func (kt *KsefTokenHandler) encryptToken(tokenPlaintext string, timestamp time.Time) (string, error) {
	certificate, err := kt.apiConfig.CertificatesDB.GetByUsage(certsdb.UsageTokenEncryption)
	if err != nil {
		return "", err
	}
	encryptedBytes, err := encryption.EncryptMessageWithCertificate(
		certificate.PEMFile,
		fmt.Appendf([]byte{}, "%s|%d", tokenPlaintext, timestamp.UnixMilli()),
	)
	if err != nil {
		return "", err
	}
	return string(encryptedBytes), nil
}
