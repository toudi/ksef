package xades

import (
	"context"
	"io"
	"ksef/internal/http"
	"ksef/internal/logging"
	"ksef/internal/sei/api/client/v2/auth/validator"
	baseHTTP "net/http"
)

const (
	endpointValidateSignedChallenge = "/api/v2/auth/xades-signature"
)

func validateSignedChallenge(ctx context.Context, httpClient *http.Client, signedChallenge io.Reader, success func(resp validator.ValidationReference)) error {
	logging.AuthLogger.Debug("validate signed challenge")

	var resp validator.ValidationReference
	var err error

	_, err = httpClient.Request(
		ctx,
		http.RequestConfig{
			ContentType:     http.XML,
			Body:            signedChallenge,
			Dest:            &resp,
			DestContentType: http.JSON,
			ExpectedStatus:  baseHTTP.StatusAccepted,
			Method:          baseHTTP.MethodPost,
		},
		endpointValidateSignedChallenge,
	)

	if err != nil {
		logging.AuthLogger.Error("error validating challenge response", "err", err)
	} else {
		logging.AuthLogger.Debug("challenge validation successful")
		go success(resp)
	}

	return err
}
