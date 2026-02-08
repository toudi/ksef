package http

import (
	ratelimiter "ksef/internal/utils/rate-limiter"
	"log/slog"
)

type RequestRateLimit struct {
	limits map[string]*ratelimiter.Limiter
	logger *slog.Logger
}

func (rl *RequestRateLimit) Wait(operationId string) {
	limitsKey := "other"

	if _, exists := rl.limits[operationId]; exists {
		limitsKey = operationId
	}

	rl.logger.Debug("using rate limits slot to determine rate limit", "key", limitsKey)

	rl.limits[limitsKey].Wait()
}

func NewRequestRateLimit(logger *slog.Logger, limits map[string]*ratelimiter.Limiter) *RequestRateLimit {
	return &RequestRateLimit{
		logger: logger,
		limits: limits,
	}
}
