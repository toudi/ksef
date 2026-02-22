package auth

import (
	"context"
	"fmt"
	"ksef/internal/http"
	baseHTTP "net/http"
)

const (
	endpointLogout = "/v2/auth/sessions/%s"
)

func (t *TokenManager) Logout(sessionRefNo string) error {
	if t.sessionTokens == nil {
		return errNoTokens
	}

	if sessionRefNo == "" {
		sessionRefNo = "current"
	}

	_, err := t.httpClient.Request(
		context.Background(),
		http.RequestConfig{
			Headers:        map[string]string{"Authorization": "Bearer " + t.sessionTokens.AuthorizationToken.Token},
			Method:         baseHTTP.MethodDelete,
			ExpectedStatus: baseHTTP.StatusNoContent,
		},

		fmt.Sprintf(endpointLogout, sessionRefNo),
	)

	return err
}
