package interactive

import (
	"embed"
	"fmt"
	"ksef/internal/encryption"
	encryptionPkg "ksef/internal/encryption"
	"ksef/internal/logging"
	"ksef/internal/sei/api/client"
	"log/slog"
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

const (
	endpointInitToken              = "/api/online/Session/InitToken"
	endpointAuthorisationChallenge = "/api/online/Session/AuthorisationChallenge"
)

type authorisationChallengeTemplateVarsType struct {
	Cipher         encryption.CipherTemplateVarsType
	Issuer         string
	Challenge      string
	EncryptedToken []byte
}

func (i *InteractiveSession) Login(issuer string, waitForAuthStatus bool) error {
	var err error
	var log *slog.Logger = logging.InteractiveLogger

	authorisationRequest.ContextIdentifier.Identifier = issuer
	if i.session == nil {
		i.session = client.NewHTTPSession(i.apiClient.Environment.Host)
	}
	encryption, err := i.apiClient.Encryption()
	if err != nil {
		return fmt.Errorf("unable to initialize encryption: %v", err)
	}

	_, err = i.session.JSONRequest(
		client.JSONRequestParams{
			Method:   "POST",
			Endpoint: endpointAuthorisationChallenge,
			Payload:  authorisationRequest,
			Response: &authorisationResponse,
			Logger:   logging.InteractiveHTTPLogger,
		},
	)

	if err != nil {
		return fmt.Errorf("unable to call authorisationRequest: %v", err)
	}

	if i.issuerToken == "" {
		gatewayToken, err := i.retrieveGateweayToken(issuer)
		if err != nil || gatewayToken == "" {
			return fmt.Errorf("cannot retrieve gateway token: %v", err)
		}
		i.issuerToken = gatewayToken
	}

	var challengePlaintext = fmt.Sprintf("%s|%d", i.issuerToken, authorisationResponse.Timestamp.UnixMilli())
	log.Debug("InteractiveSession::Login", "challengePlaintext", challengePlaintext)

	var authorisationChallengeTemplateVars = authorisationChallengeTemplateVarsType{
		Cipher:    encryption.CipherTemplateVars,
		Issuer:    issuer,
		Challenge: authorisationResponse.Challenge,
	}

	authorisationChallengeTemplateVars.EncryptedToken, err = encryptionPkg.EncryptMessageWithCertificate(
		i.apiClient.Environment.RsaPublicKey,
		[]byte(challengePlaintext),
	)
	if err != nil {
		return fmt.Errorf("error encrypting gatewayToken: %v", err)
	}
	authorisationChallengeTemplateVars.Cipher.EncryptionKey, err = encryptionPkg.EncryptMessageWithCertificate(
		i.apiClient.Environment.RsaPublicKey,
		encryption.Cipher.Key,
	)
	if err != nil {
		return fmt.Errorf("error encrypting cipher key: %v", err)
	}

	var initTokenResponse initTokenResponseType

	response, err := i.session.XMLRequest(
		client.XMLRequestParams{
			Method:       "POST",
			Endpoint:     endpointInitToken,
			TemplateDir:  authorisationChallengeTemplate,
			TemplateName: "interactive_auth_challenge.xml",
			TemplateData: authorisationChallengeTemplateVars,
			Response:     &initTokenResponse,
			Logger:       logging.InteractiveHTTPLogger,
		},
	)

	if err != nil || (response != nil && response.StatusCode/100 != 2) {
		if response != nil {
			return fmt.Errorf(
				"unexpected response code: %d != 2xx or err: %v",
				response.StatusCode,
				err,
			)
		}
		return fmt.Errorf("error processing initToken: %v", err)
	}

	i.session.SetHeader("SessionToken", initTokenResponse.Token.Value)
	i.referenceNo = initTokenResponse.ReferenceNumber

	log.Debug(
		"InteractiveSession::Login",
		"SessionToken",
		initTokenResponse.Token.Value,
		"SessionID",
		initTokenResponse.ReferenceNumber,
	)

	if waitForAuthStatus {
		return i.WaitForStatusCode(VerifySecuritySuccess, 15)
	}
	return nil
}

func init() {
	authorisationRequest.ContextIdentifier.Type = "onip"
}
