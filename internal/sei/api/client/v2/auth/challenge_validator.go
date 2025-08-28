package auth

import (
	"context"
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
	ValidateChallenge(ctx context.Context, challenge authChallengeResponse) (*ValidationResponse, error)
}
