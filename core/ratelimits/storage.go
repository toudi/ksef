// Package ratelimits defines the contract for persisting the state of the KSeF
// request rate limiter.
package ratelimits

import (
	"context"
	"time"
)

// SlotEntries holds the recorded requests that fall inside a single limiter slot
// (e.g. the "1 minute" window). It is the unit of the persisted state.
type SlotEntries struct {
	Slot    time.Duration
	Entries []time.Time
}

// State is a snapshot of the rate limiter: operation ID -> slots together with
// the requests that are still inside their sliding windows.
type State map[string][]SlotEntries

// StateStorage persists and restores the rate limiter state.
type StateStorage interface {
	// Load returns current state. a (nil, nil) is a valid response - simply means that
	// there are no limits recorded yet (e.g. when you start the client for the first time)
	Load(ctx context.Context) (State, error)
	// Set records the whole state. It is an in-memory operation: implementations are
	// expected to buffer it and write it out in Close
	Set(ctx context.Context, state State) error
	// Close persists whatever Set recorded and releases whatever the implementation
	// holds (open files, a database or cache connection pool, ...). It is called once,
	// by whoever constructed the storage.
	Close(ctx context.Context) error
}
