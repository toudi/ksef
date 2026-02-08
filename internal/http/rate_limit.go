package http

import (
	ratelimiter "ksef/internal/utils/rate-limiter"
	"log/slog"
)

type RequestRateLimit struct {
	limits map[string]*ratelimiter.Limiter
	logger *slog.Logger
}

func (rl *RequestRateLimit) limitsKey(operationId string) string {
	limitsKey := "other"

	if _, exists := rl.limits[operationId]; exists {
		limitsKey = operationId
	}

	return limitsKey
}

func (rl *RequestRateLimit) Wait(operationId string) {
	limitsKey := rl.limitsKey(operationId)
	rl.logger.Debug("using rate limits slot to determine rate limit", "key", limitsKey)
	rl.limits[limitsKey].Wait()
}

// this function is essentially a hacky way of rewriting time that last request was made
// by default, we record our client's time - but we don't know how does the ministry's server
// record it and so it is safer to replace the entry by inspecting the client's time (again)
// after the response comes back to us.
func (rl *RequestRateLimit) replaceLastEntry(operationId string) {
	limitsKey := rl.limitsKey(operationId)
	rl.logger.Debug("replace last entry in limits slot", "key", limitsKey)
	rl.limits[limitsKey].ReplaceLastEntry()
}

func NewRequestRateLimit(logger *slog.Logger, limits map[string]*ratelimiter.Limiter) *RequestRateLimit {
	return &RequestRateLimit{
		logger: logger,
		limits: limits,
	}
}
