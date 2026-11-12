package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"

	"ksef/core/certsdb"
	"ksef/core/keyring"
	"ksef/internal/logging"
)

var (
	errEmptySessionTokens = errors.New("session tokens are empty")
	errCannotUseToken     = errors.New("token cannot be used")
)

func (tm *TokenManager) persistTokens() error {
	logger := logging.AuthLogger.With("auth", "token manager")
	logger.Debug("zapisywanie tokenów sesyjnych")

	// Snapshot the tokens under the lock, then release it before doing I/O.
	tm.mutex.Lock()
	sessionTokens := tm.sessionTokens
	tm.mutex.Unlock()

	if sessionTokens == nil {
		logger.Error("tokeny sesyjne nieustawione")
		return errEmptySessionTokens
	}
	var buffer bytes.Buffer
	if err := json.NewEncoder(&buffer).Encode(sessionTokens); err != nil {
		logger.Error("enkodowanie tokenów do JSON zakończone niepowodzeniem", "err", err)
		return err
	}
	authCert, err := tm.certsDB.GetCertificateByUsageAndNIP(
		context.Background(), certsdb.CertificateUsageAuthentication, tm.nip,
	)
	if err != nil {
		logger.Error("nie udało się odczytać certyfikatu", "err", err)
		return err
	}

	err = tm.keyring.Set(tm.environmentID, tm.nip, keyring.SessionTokensKey(authCert.UID), buffer.String())

	logger.Debug("rezultat zapisywania tokenów do keyringu", "err", err)
	return err
}

func (tm *TokenManager) clearSessionTokens() error {
	logger := logging.AuthLogger.With("auth", "token manager")

	tm.clearAuthorizationToken()
	authCert, err := tm.certsDB.GetCertificateByUsageAndNIP(
		context.Background(), certsdb.CertificateUsageAuthentication, tm.nip,
	)
	if err != nil {
		logger.Error("nie udało się odczytać certyfikatu", "err", err)
		return err
	}

	return tm.keyring.Delete(tm.environmentID, tm.nip, keyring.SessionTokensKey(authCert.UID))
}

func (tm *TokenManager) restoreTokens(ctx context.Context) error {
	logger := logging.AuthLogger.With("auth", "token manager")
	logger.Debug("próba odczytania tokenów z systemowego pęku kluczy")
	logger.Debug("próba odczytania certyfikatu odpowiedzialnego za autoryzację")
	authCert, err := tm.certsDB.GetCertificateByUsageAndNIP(
		ctx, certsdb.CertificateUsageAuthentication, tm.nip,
	)
	if err != nil {
		logger.Error("nie udało się odczytać certyfikatu", "err", err)
		return err
	}
	tokens, err := tm.keyring.Get(tm.environmentID, tm.nip, keyring.SessionTokensKey(authCert.UID))
	if err != nil && err != keyring.ErrNotFound && err != keyring.ErrPermissionDenied {
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
		tm.updateSessionTokens(&sessionTokens, refreshed)
		tm.markReady()
		return nil
	}
	logger.Debug("tokeny nie mogą być użyte ponownie - usuwam z pęku kluczy")
	if err = tm.clearSessionTokens(); err != nil {
		logger.Error("błąd czyszczenia tokenów", "err", err)
		return err
	}
	return errCannotUseToken
}
