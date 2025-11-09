package auth

import (
	"context"
	"errors"
	"fmt"
	"ksef/internal/http"
	baseHTTP "net/http"
)

var (
	errNoTokens = errors.New("no tokens")
)

const (
	endpointLogout = "/api/v2/auth/sessions/%s"
)

func (t *TokenManager) Logout() error {
	if t.sessionTokens == nil {
		return errNoTokens
	}

	_, err := t.httpClient.Request(
		context.Background(),
		http.RequestConfig{
			Headers:        map[string]string{"Authorization": "Bearer " + t.sessionTokens.AuthorizationToken.Token},
			Method:         baseHTTP.MethodDelete,
			ExpectedStatus: baseHTTP.StatusNoContent,
		},
		fmt.Sprintf(endpointLogout, t.validationReference.ReferenceNumber),
	)

	return err
}
