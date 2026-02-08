package ratelimiter

import (
	"ksef/internal/logging"
	"time"
)

type Limiter struct {
	limits []RateLimit
}

func (l *Limiter) Wait() {
	var allowed bool = false
	var waitTime time.Duration
	var now time.Time

	for !allowed {
		now = time.Now()

		for limitIdx := range l.limits {
			if l.limits[limitIdx].buffer == nil {
				l.limits[limitIdx].buffer = NewCallHistory(l.limits[limitIdx].Limit)
			}
			limiter := l.limits[limitIdx]
			allowed, waitTime = limiter.buffer.IsAllowed(limiter.Slot, now)
			if !allowed {
				break
			}
		}
		if waitTime > 0 {
			logging.HTTPLogger.Debug("Rate limit exceeded; waiting for next slot", "sleep", waitTime)
			time.Sleep(waitTime)
		}
	}

	// request is allowed - add entry to rate limits.
	for limitIdx := range l.limits {
		l.limits[limitIdx].buffer.Push(now)
	}
}

func NewLimiter(limits []RateLimit) *Limiter {
	return &Limiter{
		limits: limits,
	}
}
