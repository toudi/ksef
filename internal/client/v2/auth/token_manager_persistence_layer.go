package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"ksef/internal/certsdb"
	"ksef/internal/keyring"
	"ksef/internal/logging"
	"ksef/internal/runtime"
)

var (
	errEmptySessionTokens = errors.New("session tokens are empty")
	errCannotUseToken     = errors.New("token cannot be used")
)

func (tm *TokenManager) persistTokens() error {
	logger := logging.AuthLogger.With("auth", "token manager")
	logger.Debug("zapisywanie tokenów sesyjnych")
	if tm.sessionTokens == nil {
		logger.Error("tokeny sesyjne nieustawione")
		return errEmptySessionTokens
	}
	var buffer bytes.Buffer
	if err := json.NewEncoder(&buffer).Encode(tm.sessionTokens); err != nil {
		logger.Error("enkodowanie tokenów do JSON zakończone niepowodzeniem", "err", err)
		return err
	}
	nip, err := runtime.GetNIP(tm.vip)
	if err != nil {
		logger.Error("nieprawidłowy numer NIP", "err", err)
		return err
	}
	authCert, err := tm.certsDB.GetByUsage(certsdb.UsageAuthentication, nip)
	if err != nil {
		logger.Error("nie udało się odczytać certyfikatu", "err", err)
	}

	err = tm.keyring.Set(string(runtime.GetEnvironmentId(tm.vip)), nip, keyring.SessionTokensKey(authCert.UID), buffer.String())

	logger.Debug("rezultat zapisywania tokenów do keyringu", "err", err)
	return err
}

func (tm *TokenManager) SetSessionTokens(tokens *SessionTokens) {
	tm.sessionTokens = tokens
}

func (tm *TokenManager) clearSessionTokens() error {
	logger := logging.AuthLogger.With("auth", "token manager")

	tm.sessionTokens = nil
	nip, err := runtime.GetNIP(tm.vip)
	if err != nil {
		return err
	}
	authCert, err := tm.certsDB.GetByUsage(certsdb.UsageAuthentication, nip)
	if err != nil {
		logger.Error("nie udało się odczytać certyfikatu", "err", err)
	}

	return tm.keyring.Delete(string(runtime.GetEnvironmentId(tm.vip)), nip, keyring.SessionTokensKey(authCert.UID))
}

func (tm *TokenManager) restoreTokens(ctx context.Context) error {
	logger := logging.AuthLogger.With("auth", "token manager")
	logger.Debug("próba odczytania tokenów z systemowego pęku kluczy")
	nip, err := runtime.GetNIP(tm.vip)
	if err != nil {
		logger.Error("błąd odczytu numeru NIP", "err", err)
		return err
	}
	gateway := runtime.GetEnvironmentId(tm.vip)
	logger.Debug("próba odczytania certyfikatu odpowiedzialnego za autoryzację")
	authCert, err := tm.certsDB.GetByUsage(certsdb.UsageAuthentication, nip)
	if err != nil {
		logger.Error("nie udało się odczytać certyfikatu", "err", err)
	}
	tokens, err := tm.keyring.Get(string(gateway), nip, keyring.SessionTokensKey(authCert.UID))
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
	canBeUsed, refreshed := tm.validateSessionTokens(ctx, &sessionTokens)
	if canBeUsed {
		logger.Debug("tokeny mogą być użyte ponownie")
		tm.sessionTokens = &sessionTokens
		if refreshed {
			if err = tm.persistTokens(); err != nil {
				logger.Error("błąd zapisu odświeżonych tokenów", "err", err)
				return err
			}
		}
		if tm.vip.GetBool(FlagExitAfterPersistingToken) {
			tm.finished = true
		}
		go func() {
			tm.updateChannel <- TokenUpdate{
				Token: sessionTokens.AuthorizationToken.Token,
			}
		}()
		return nil
	}
	logger.Debug("tokeny nie mogą być użyte ponownie - usuwam z pęku kluczy")
	if err = tm.clearSessionTokens(); err != nil {
		logger.Error("błąd czyszczenia tokenów", "err", err)
		return err
	}
	return errCannotUseToken
}
