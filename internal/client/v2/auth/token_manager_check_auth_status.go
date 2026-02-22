package auth

import (
	"context"
	"fmt"
	"ksef/internal/http"
	"ksef/internal/logging"
	baseHttp "net/http"
)

const (
	endpointAuthStatus        = "/v2/auth/%s"
	authStatusCodeSuccess int = 200
)

func (t *TokenManager) checkAuthStatus(ctx context.Context) error {
	authHeaders := map[string]string{
		"Authorization": "Bearer " + t.validationReference.AuthenticationToken.Token,
	}

	var authStatus AuthenticationStatus

	_, err := t.httpClient.Request(ctx, http.RequestConfig{
		Headers:         authHeaders,
		Dest:            &authStatus,
		DestContentType: http.JSON,
		ExpectedStatus:  baseHttp.StatusOK,
		Method:          baseHttp.MethodGet,
	}, fmt.Sprintf(endpointAuthStatus, t.validationReference.ReferenceNumber))

	logging.AuthLogger.Debug("auth response", "content", authStatus)

	if err != nil {
		logging.AuthLogger.Error("error checking auth status", "err", err)
		return err
	}

	if authStatus.Status.Code == authStatusCodeSuccess {
		logging.AuthLogger.Debug("authentication is successful - proceed to redeeming tokens")
		return t.redeemTokens(ctx)
	}

	return nil
}
