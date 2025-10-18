package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"ksef/internal/environment"

	"github.com/zalando/go-keyring"
)

var (
	errEmptySessionTokens = errors.New("session tokens are empty")
)

func (tm *TokenManager) PersistTokens(env environment.Environment, nip string) error {
	if tm.sessionTokens == nil {
		return errEmptySessionTokens
	}
	var buffer bytes.Buffer
	if err := json.NewEncoder(&buffer).Encode(tm.sessionTokens); err != nil {
		return err
	}
	return keyring.Set(string(env)+"-sessionTokens", nip, buffer.String())
}

func (tm *TokenManager) SetSessionTokens(tokens *SessionTokens) {
	tm.sessionTokens = tokens
}
