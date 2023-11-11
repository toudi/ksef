package api

import (
	"embed"
	"fmt"
	"ksef/common"
	"time"
)

type authorisationRequestType struct {
	ContextIdentifier struct {
		Type       string `json:"type"`
		Identifier string `json:"identifier"`
	} `json:"contextIdentifier"`
}

type authorisationResponseType struct {
	Timestamp time.Time `json:"timestamp"`
	Challenge string    `json:"challenge"`
}

type initTokenResponseType struct {
	ReferenceNumber string `json:"referenceNumber"`
	Token           struct {
		Value string `json:"token"`
	} `json:"sessionToken"`
}

var authorisationRequest authorisationRequestType
var authorisationResponse authorisationResponseType

//go:embed "interactive_auth_challenge.xml"
var authorisationChallengeTemplate embed.FS

type authorisationChallengeTemplateVarsType struct {
	Cipher         cipherTemplateVarsType
	Issuer         string
	Challenge      string
	EncryptedToken []byte
}

func (i *InteractiveSession) login(issuer string) error {
	var err error

	authorisationRequest.ContextIdentifier.Identifier = issuer

	_, err = i.api.requestFactory.JSONRequest("POST", EndpointAuthorisationChallenge, authorisationRequest, &authorisationResponse)

	if err != nil {
		return fmt.Errorf("unable to call authorisationRequest: %v", err)
	}

	gatewayToken, err := i.retrieveGateweayToken(issuer)
	if err != nil || gatewayToken == "" {
		return fmt.Errorf("cannot retrieve gateway token: %v", err)
	}

	var challengePlaintext = fmt.Sprintf("%s|%d", gatewayToken, authorisationResponse.Timestamp.UnixMilli())
	var authorisationChallengeTemplateVars = authorisationChallengeTemplateVarsType{
		Cipher:    i.api.cipherTemplateVars,
		Issuer:    issuer,
		Challenge: authorisationResponse.Challenge,
	}

	authorisationChallengeTemplateVars.EncryptedToken, err = common.EncryptMessageWithCertificate(i.api.environment.rsaPublicKey, []byte(challengePlaintext))
	if err != nil {
		return fmt.Errorf("error encrypting gatewayToken: %v", err)
	}
	authorisationChallengeTemplateVars.Cipher.EncryptionKey, err = common.EncryptMessageWithCertificate(i.api.environment.rsaPublicKey, i.api.cipher.Key)
	if err != nil {
		return fmt.Errorf("error encrypting cipher key: %v", err)
	}

	var initTokenResponse initTokenResponseType

	response, err := i.api.requestFactory.XMLRequest(
		"POST",
		EndpointInitToken,
		authorisationChallengeTemplate,
		"interactive_auth_challenge.xml",
		authorisationChallengeTemplateVars,
		&initTokenResponse,
	)

	if err != nil || (response != nil && response.StatusCode/100 != 2) {
		if response != nil {
			return fmt.Errorf("unexpected response code: %d != 2xx or err: %v", response.StatusCode, err)
		}
		return fmt.Errorf("error processing initToken: %v", err)
	}

	i.api.requestFactory.headers.Add("SessionToken", initTokenResponse.Token.Value)
	i.referenceNo = initTokenResponse.ReferenceNumber

	return nil
}

func init() {
	authorisationRequest.ContextIdentifier.Type = "onip"
}
