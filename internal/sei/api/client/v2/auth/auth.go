package auth

import (
	"errors"
	"ksef/internal/config"
	"ksef/internal/http"
)

var ErrInvalidChallengeValidator = errors.New("invalid challenge validator")

type AuthHandler struct {
	httpClient         http.Client
	config             config.APIConfig
	challengeValidator AuthValidator
	tokenManager       TokenManager
}

type Manager struct {
	handler *AuthHandler
}

func NewAuthHandler(httpClient http.Client, config config.APIConfig, opts ...func(*AuthHandler)) *AuthHandler {
	handler := &AuthHandler{
		httpClient: httpClient,
		config:     config,
	}

	for _, opt := range opts {
		opt(handler)
	}

	return handler
}

func NewManager(handler *AuthHandler) *Manager {
	return &Manager{handler: handler}
}

func (a *AuthHandler) Close() {
	a.tokenManager.Stop()
}

func WithChallengeValidator(validator AuthValidator) func(handler *AuthHandler) {
	return func(handler *AuthHandler) {
		handler.challengeValidator = validator
	}
}

func (am *Manager) SetKsefToken(token string) error {
	// this only works for challenge validator set to ksef so let's assert that
	ksefTokenValidator, ok := am.handler.challengeValidator.(*ksefTokenAuthValidator)
	if !ok {
		return ErrInvalidChallengeValidator
	}
	ksefTokenValidator.ksefToken = token
	return nil
}
