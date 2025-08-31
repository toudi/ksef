package auth

import (
	"context"
	"ksef/internal/http"
	"time"
)

type ValidationAuthenticationToken struct {
	Token      string    `json:"token"`
	ValidUntil time.Time `json:"validUntil"`
}

type ValidationResponse struct {
	ReferenceNumber     string                        `json:"referenceNumber"`
	AuthenticationToken ValidationAuthenticationToken `json:"authenticationToken"`
}

type AuthValidator interface {
	ValidateChallenge(ctx context.Context, httpClient http.Client, challenge authChallengeResponse) (*ValidationResponse, error)
}
