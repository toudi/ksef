package esignature

import (
	"context"
	"fmt"
	"ksef/internal/http"
	"ksef/internal/sei/api/client/v2/auth/validator"
	baseHTTP "net/http"
	"os"
)

const (
	endpointValidateSignedChallenge = "/api/v2/auth/xades-signature"
)

func (e *eSignatureTokenHandler) validateSignedChallenge(signedFilePath string) error {
	var resp validator.ValidationReference
	var ctx = context.Background()

	signedFile, err := os.Open(signedFilePath)
	if err != nil {
		return err
	}
	defer signedFile.Close()

	_, err = e.httpClient.Request(
		ctx,
		http.RequestConfig{
			ContentType:     http.XML,
			Body:            signedFile,
			Dest:            &resp,
			DestContentType: http.JSON,
			ExpectedStatus:  baseHTTP.StatusAccepted,
			Method:          baseHTTP.MethodPost,
		},
		endpointValidateSignedChallenge,
	)

	fmt.Printf("challenge validation response: %+v; err=%v\n", resp, err)

	if err == nil {
		go func() {
			e.eventChannel <- validator.AuthEvent{
				State:               validator.StateValidationReferenceResult,
				ValidationReference: &resp,
			}
		}()
	}

	return err
}
