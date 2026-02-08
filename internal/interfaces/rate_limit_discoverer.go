package interfaces

import (
	"context"
	ratelimiter "ksef/internal/utils/rate-limiter"
)

type RateLimitsDiscoverFunc func(ctx context.Context, host string, authToken string) (map[string]*ratelimiter.Limiter, error)
