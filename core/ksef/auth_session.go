package ksef

import "time"

// AuthSessionStatus is the status block of a KSeF authentication session.
type AuthSessionStatus struct {
	Code        int    `json:"status"`
	Description string `json:"description"`
}

// AuthSession is a single KSeF authentication session as returned by
// GET /v2/auth/sessions.
type AuthSession struct {
	ID                     string            `json:"referenceNumber"`
	Current                bool              `json:"isCurrent"`
	StartDate              time.Time         `json:"startDate"`
	AuthMethod             string            `json:"authenticationMethod"`
	Status                 AuthSessionStatus `json:"status"`
	TokenRedeemed          bool              `json:"isTokenRedeemed"`
	RefreshTokenValidUntil time.Time         `json:"refreshTokenValidUntil"`
}

// AuthSessions is the paginated response of GET /v2/auth/sessions.
type AuthSessions struct {
	Sessions []AuthSession `json:"items"`
}
