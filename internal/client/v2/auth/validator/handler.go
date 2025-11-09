package validator

import (
	"context"
	"ksef/internal/http"
	"time"
)

type AuthChallenge struct {
	Challenge string    `json:"challenge"`
	Timestamp time.Time `json:"timestamp"`
}

type ValidationAuthenticationToken struct {
	Token      string    `json:"token"`
	ValidUntil time.Time `json:"validUntil"`
}

type ValidationReference struct {
	ReferenceNumber     string                        `json:"referenceNumber"`
	AuthenticationToken ValidationAuthenticationToken `json:"authenticationToken"`
}

type State uint8

const (
	StateInitialized                 State = iota
	StateAwaitingChallenge           State = iota
	StateAwaitingChallengeValidation State = iota
	StateValidationReferenceResult   State = iota
	// when tokens are retrieved from the safe storage
	StateTokensRestored State = iota
	StateExit           State = iota
)

type AuthEvent struct {
	// when EventType == EventTypeInitialized, then it means that validator has all the
	// information requried to validate a challenge and the token manager can start
	// auth challenge procedure
	State               State
	ValidationReference *ValidationReference
	SessionTokens       string
	Error               error
}

type AuthChallengeValidator interface {
	Event() chan AuthEvent
	Initialize(ctx context.Context, httpClient *http.Client) error
	ValidateChallenge(ctx context.Context, challenge AuthChallenge) error
}
