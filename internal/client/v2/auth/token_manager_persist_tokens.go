package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"ksef/internal/config"
	"ksef/internal/logging"

	"github.com/zalando/go-keyring"
)

var (
	errEmptySessionTokens = errors.New("session tokens are empty")
)

func (tm *TokenManager) PersistTokens(gw config.Gateway, nip string) error {
	logging.AuthLogger.Debug("zapisywanie tokenów sesyjnych")
	if tm.sessionTokens == nil {
		return errEmptySessionTokens
	}
	var buffer bytes.Buffer
	if err := json.NewEncoder(&buffer).Encode(tm.sessionTokens); err != nil {
		return err
	}
	return keyring.Set(string(gw)+"-sessionTokens", nip, buffer.String())
}

func (tm *TokenManager) SetSessionTokens(tokens *SessionTokens) {
	tm.sessionTokens = tokens
}
