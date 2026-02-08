package ratelimiter

import (
	"fmt"
	"time"
)

// for each of the limited endpoints, we need to keep track of
// the previous HTTP requests. The server uses a sliding window
// technique to limit the number of requests therefore we need
// to implement a similar thing in order to know ahead of time
// if we can make a query or not.
type callHistory struct {
	entries    []time.Time
	capacity   int
	writeIndex int
}

func NewCallHistory(capacity int) *callHistory {
	return &callHistory{
		entries:  make([]time.Time, 0, capacity),
		capacity: capacity,
	}
}

func (b *callHistory) Push(value time.Time) {
	if len(b.entries) < b.capacity {
		b.entries = append(b.entries, value)
	} else {
		b.entries[b.writeIndex] = value
		b.writeIndex = (b.writeIndex + 1) % len(b.entries)
	}
}

// this function checks if there are less entries than now - slot than capacity.
// if not - it returns the time difference betwen earliest entry and now that needs
// to be waited for the slot to become available
func (b *callHistory) IsAllowed(slot time.Duration, now time.Time) (isAllowed bool, waitTime time.Duration) {
	// if we can skip lookups - let's do that
	if len(b.entries) < b.capacity {
		return true, 0
	}

	// ok we're here so let's start by picking up the smallest entry
	var earliestEntry time.Time
	// so because the rate limit uses sort of a "sliding window" (? don't know how to call it better)
	// we need to keep track of number of entries within the slot / window.
	// for example:
	// if a "slot" is 1 minute then we need to inspect all of the entries that are after now - 1 minute.
	cutoff := now.Add(-slot)

	// how many entries are occuppied within the slot / sliding window
	var occupied int

	for _, entry := range b.entries {
		if entry.After(cutoff) {
			occupied += 1
			if entry.Before(earliestEntry) || earliestEntry.IsZero() {
				earliestEntry = entry
			}
		}
	}

	// perfect. now we can determine if we're over the limit:
	isAllowed = occupied < b.capacity
	// and, if not - how much do we have to wait
	if !isAllowed {
		waitTime = earliestEntry.Add(slot).Sub(now)
		fmt.Printf("calculated wait time as %v\n", waitTime)
	}
	return isAllowed, waitTime
}
