package http

import (
	"errors"
	"log/slog"
	"time"

	"ksef/core/ratelimits"
	"ksef/internal/logging"

	ratelimiter "ksef/internal/utils/rate-limiter"
)

var errLimitNotFound = errors.New("limit not found")

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

func (rl *RequestRateLimit) Wait(operationId string) time.Duration {
	limitsKey := rl.limitsKey(operationId)
	rl.logger.Debug("using rate limits slot to determine rate limit", "key", limitsKey)
	return rl.limits[limitsKey].Wait(func(sleepFor, remaining time.Duration) {
		rl.logger.Debug("rate limit exceeded", "sleepFor", sleepFor, "remaining", remaining)
	})
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

func (rl *RequestRateLimit) GetLimit(operationId string, slot time.Duration) (int, error) {
	limitsKey := rl.limitsKey(operationId)
	limiter, exists := rl.limits[limitsKey]
	if !exists {
		return 0, errLimitNotFound
	}

	for _, limit := range limiter.Limits() {
		if limit.Slot == slot {
			return limit.Limit, nil
		}
	}

	return 0, errLimitNotFound
}

// Snapshot converts the limiter into its persistable form: for every operation ID
// it holds the requests that are still within the limiter's sliding windows.
func (rl *RequestRateLimit) snapshot(now time.Time) ratelimits.State {
	state := make(ratelimits.State, len(rl.limits))

	for operationId, limiter := range rl.limits {
		requestsWithinSlot := limiter.EntriesWithinSlot(now)
		if len(requestsWithinSlot) == 0 {
			continue
		}

		slots := make([]ratelimits.SlotEntries, len(requestsWithinSlot))
		for idx, entry := range requestsWithinSlot {
			slots[idx] = ratelimits.SlotEntries{
				Slot:    entry.Slot,
				Entries: entry.Entries,
			}
		}
		state[operationId] = slots
	}

	return state
}

// Load restores the limiter from a previously persisted state. Entries for
// unknown operation IDs are ignored.
func (rl *RequestRateLimit) load(state ratelimits.State) {
	for operationId, slots := range state {
		limiter, exists := rl.limits[operationId]
		if !exists {
			logging.HTTPLogger.Warn("rate limiter for operation ID does not exist", "operationId", operationId)
			continue
		}

		entries := make([]ratelimiter.LimiterSlotEntries, len(slots))
		for idx, slot := range slots {
			entries[idx] = ratelimiter.LimiterSlotEntries{
				Slot:    slot.Slot,
				Entries: slot.Entries,
			}
		}
		limiter.LoadEntries(entries)
	}
}
