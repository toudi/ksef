package auth

import (
	"ksef/internal/config"
	"ksef/internal/http"
)

type AuthHandler struct {
	httpClient         http.Client
	config             config.APIConfig
	challengeValidator AuthValidator
	tokenManager       TokenManager
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

func (a *AuthHandler) Close() {
	a.tokenManager.Stop()
}

func WithChallengeValidator(validator AuthValidator) func(handler *AuthHandler) {
	return func(handler *AuthHandler) {
		handler.challengeValidator = validator
	}
}
