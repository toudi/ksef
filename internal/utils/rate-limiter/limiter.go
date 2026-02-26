package ratelimiter

import (
	"time"

	"github.com/samber/lo"
)

type Limiter struct {
	limits []RateLimit
}

func (l *Limiter) Wait(progress func(sleepTime, remaining time.Duration)) time.Duration {
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
			sleepWithProgressFunc(waitTime, progress)
		}
	}

	// request is allowed - add entry to rate limits.
	for limitIdx := range l.limits {
		l.limits[limitIdx].buffer.Push(now)
	}

	return waitTime
}

func NewLimiter(limits []RateLimit) *Limiter {
	return &Limiter{
		limits: limits,
	}
}

type LimiterSlotEntries struct {
	Slot    time.Duration `yaml:"slot"`
	Entries []time.Time   `yaml:"requests"`
}

func (l *Limiter) EntriesWithinSlot(now time.Time) (result []LimiterSlotEntries) {
	for limitIdx := range l.limits {
		limiter := l.limits[limitIdx]
		if limiter.buffer == nil {
			continue
		}

		// there is no point in persisting expired requests so let's save only these which
		// are within the sliding window
		minTime := now.Add(-limiter.Slot)
		entriesWithinSlot := lo.Filter(limiter.buffer.entries, func(entry time.Time, _ int) bool {
			return entry.After(minTime)
		})

		if len(entriesWithinSlot) > 0 {
			result = append(result, LimiterSlotEntries{
				Slot:    limiter.Slot,
				Entries: entriesWithinSlot,
			})
		}

	}

	return result
}

func (l *Limiter) LoadEntries(entries []LimiterSlotEntries) {
	for _, slottedEntries := range entries {
		// we've got to lookup the limiter by slot
		// suboptimal but there are only 3 limiters so ..
		for limiterIdx := range l.limits {
			if l.limits[limiterIdx].Slot == slottedEntries.Slot {
				// initialize the buffer (it is being lazily-loaded within the Wait() function by default
				// but here we're loading values therefore we have to initialize it.
				l.limits[limiterIdx].buffer = NewCallHistory(l.limits[limiterIdx].Limit)
				// because we have initialized the buffer above, it is safe to call append - go will not reallocate memory
				l.limits[limiterIdx].buffer.entries = append(l.limits[limiterIdx].buffer.entries, slottedEntries.Entries...)
				break
			}
		}
	}
}

func (l *Limiter) ReplaceLastEntry() {
	now := time.Now()
	for limitIdx := range l.limits {
		limiter := l.limits[limitIdx]
		if limiter.buffer == nil {
			continue
		}
		limiter.buffer.ReplaceLastEntry(now)
	}
}
