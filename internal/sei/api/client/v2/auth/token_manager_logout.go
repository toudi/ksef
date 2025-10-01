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
			Method:         baseHTTP.MethodDelete,
			ExpectedStatus: baseHTTP.StatusGone,
		},
		fmt.Sprintf(endpointLogout, t.validationReference.ReferenceNumber),
	)

	return err
}
