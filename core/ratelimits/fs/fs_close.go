package fs

import (
	"context"
	"io"
	"ksef/internal/utils"
)

// Close persists the in-memory state (if it was modified since the last Close) and
// releases the storage. The write goes through a temporary file which is renamed
// over the state file, so an interrupted save cannot leave a truncated state behind.
func (s *FSStorage) Close(ctx context.Context) error {
	if !s.dirty {
		return nil
	}
	if err := ctx.Err(); err != nil {
		return err
	}

	onDisk := make(fsState, len(s.state))
	for operationId, entries := range s.state {
		slottedEntries := make([]fsSlotEntries, len(entries))
		for idx, entry := range entries {
			slottedEntries[idx] = fsSlotEntries{
				Slot:    entry.Slot,
				Entries: entry.Entries,
			}
		}
		onDisk[operationId] = slottedEntries
	}

	if err := utils.WriteFileAtomic(s.dbFile, func(w io.Writer) error {
		return utils.WriteYAML(w, onDisk)
	}); err != nil {
		return err
	}

	s.dirty = false

	return nil
}
