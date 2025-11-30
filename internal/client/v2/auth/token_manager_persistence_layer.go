package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"ksef/internal/logging"
	"ksef/internal/runtime"

	"github.com/zalando/go-keyring"
)

var (
	errEmptySessionTokens = errors.New("session tokens are empty")
	errCannotUseToken     = errors.New("token cannot be used")
)

func keyringKey(gw runtime.Gateway) string {
	return string(gw) + "-sessionTokens"
}

func (tm *TokenManager) persistTokens() error {
	logging.AuthLogger.Debug("zapisywanie tokenów sesyjnych")
	if tm.sessionTokens == nil {
		logging.AuthLogger.Error("tokeny sesyjne nieustawione")
		return errEmptySessionTokens
	}
	var buffer bytes.Buffer
	if err := json.NewEncoder(&buffer).Encode(tm.sessionTokens); err != nil {
		logging.AuthLogger.Error("enkodowanie tokenów do JSON zakończone niepowodzeniem", "err", err)
		return err
	}
	nip, err := runtime.GetNIP(tm.vip)
	if err != nil {
		logging.AuthLogger.Error("nieprawidłowy numer NIP", "err", err)
		return err
	}
	err = keyring.Set(keyringKey(runtime.GetGateway(tm.vip)), nip, buffer.String())
	logging.AuthLogger.Debug("rezultat zapisywania tokenów do keyringu", "err", err)
	return err
}

func (tm *TokenManager) SetSessionTokens(tokens *SessionTokens) {
	tm.sessionTokens = tokens
}

func (tm *TokenManager) clearSessionTokens() error {
	tm.sessionTokens = nil
	nip, err := runtime.GetNIP(tm.vip)
	if err != nil {
		return err
	}
	return keyring.Delete(keyringKey(runtime.GetGateway(tm.vip)), nip)
}

func (tm *TokenManager) restoreTokens(ctx context.Context) error {
	var logger = logging.AuthLogger.With("auth", "token manager")
	logger.Debug("próba odczytania tokenów z systemowego pęku kluczy")
	nip, err := runtime.GetNIP(tm.vip)
	if err != nil {
		logger.Error("błąd odczytu numeru NIP", "err", err)
		return err
	}
	gateway := runtime.GetGateway(tm.vip)
	tokens, err := keyring.Get(keyringKey(gateway), nip)
	if err != nil && err != keyring.ErrNotFound {
		logger.Error("błąd odczytu tokenów", "err", err)
		return err
	}
	if tokens == "" {
		logger.Debug("brak zapisanych tokenów")
		return errCannotUseToken
	}
	var buffer bytes.Buffer
	if _, err := buffer.WriteString(tokens); err != nil {
		logger.Error("unable to write tokens to buffer")
	}
	var sessionTokens SessionTokens
	if err := json.NewDecoder(&buffer).Decode(&sessionTokens); err != nil {
		logger.Error("unable to decode tokens")
	}
	// because we're restoring tokens, let's check if we can use them
	// maybe the client already loggout out previously and so the token is invalid ?
	if tm.validateSessionTokens(ctx, &sessionTokens) {
		logger.Debug("tokeny mogą być użyte ponownie")
		tm.sessionTokens = &sessionTokens
		if err = tm.persistTokens(); err != nil {
			logger.Error("błąd zapisu odświeżonych tokenów", "err", err)
			return err
		}
		return nil
	}
	logger.Debug("tokeny nie mogą być użyte ponownie - usuwam z pęku kluczy")
	if err = tm.clearSessionTokens(); err != nil {
		logger.Error("błąd czyszczenia tokenów", "err", err)
		return err
	}
	return errCannotUseToken
}
