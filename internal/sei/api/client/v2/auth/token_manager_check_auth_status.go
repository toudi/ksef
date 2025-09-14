package auth

import (
	"context"
	"fmt"
	"ksef/internal/http"
	baseHttp "net/http"
)

const (
	endpointAuthStatus    = "/api/v2/auth/%s"
	authStatusCodeSuccess = 200
)

func (t *TokenManager) checkAuthStatus(ctx context.Context) error {
	var authHeaders = map[string]string{
		"Authorization": "Bearer " + t.validationReference.AuthenticationToken.Token,
	}

	var authStatus AuthenticationStatus

	_, err := t.httpClient.Request(ctx, http.RequestConfig{
		Headers:         authHeaders,
		Dest:            &authStatus,
		DestContentType: http.JSON,
		ExpectedStatus:  baseHttp.StatusOK,
		Method:          baseHttp.MethodPost,
	}, fmt.Sprintf(endpointAuthStatus, t.validationReference.ReferenceNumber))

	if err != nil {
		return err
	}

	if authStatus.Status.Code == authStatusCodeSuccess {
		return t.redeemTokens(ctx)
	}

	return nil
}
