package ratelimiter

import "time"

type RateLimit struct {
	Slot   time.Duration
	Limit  int
	buffer *callHistory
}
