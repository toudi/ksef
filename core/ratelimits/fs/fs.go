// This package implements the filesystem-based storage of rate-limits.
// This is done so that the internal http library can be viper-free (the previous implementation hard-wired it
// which would make it hard to make the code portable and reusable as a library)
package fs

import (
	"context"
	"ksef/core/ratelimits"
	"path/filepath"
	"time"
)

const stateFilename = "rate-limits-state.yaml"

var _ ratelimits.StateStorage = (*FSStorage)(nil)

type FSStorage struct {
	dbFile string

	state ratelimits.State
	dirty bool
}

type fsSlotEntries struct {
	Slot    time.Duration `yaml:"slot"`
	Entries []time.Time   `yaml:"requests"`
}

// fsState is the on-disk representation of ratelimits.State.
type fsState map[string][]fsSlotEntries

// NewFSStorage returns a storage that keeps the state under root and reads any
// already persisted state. A missing file is not an error.
func NewFSStorage(root string) *FSStorage {
	return &FSStorage{
		dbFile: filepath.Join(root, stateFilename),
	}
}

// Load returns the state read from disk when the storage was constructed, updated by
// any subsequent Set calls.
func (s *FSStorage) Load(ctx context.Context) (ratelimits.State, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	state, err := s.read()
	if err != nil {
		return nil, err
	}
	s.state = state

	return s.state, nil
}

// Set replaces the in-memory state. Nothing is written to disk until Close is
// called.
func (s *FSStorage) Set(ctx context.Context, state ratelimits.State) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	s.state = state
	s.dirty = true

	return nil
}
