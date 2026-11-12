package auth

import (
	"context"
	"fmt"
	baseHTTP "net/http"

	"ksef/internal/http"
)

const (
	endpointLogout = "/v2/auth/sessions/%s"
)

func (t *TokenManager) Logout(sessionRefNo string) error {
	sessionTokens := t.currentSessionTokens()
	if sessionTokens == nil || sessionTokens.AuthorizationToken == nil {
		return errNoTokens
	}

	if sessionRefNo == "" {
		sessionRefNo = "current"
	}

	_, err := t.httpClient.Request(
		context.Background(),
		http.RequestConfig{
			Headers:        map[string]string{"Authorization": "Bearer " + sessionTokens.AuthorizationToken.Token},
			Method:         baseHTTP.MethodDelete,
			ExpectedStatus: baseHTTP.StatusNoContent,
		},

		fmt.Sprintf(endpointLogout, sessionRefNo),
	)

	return err
}
