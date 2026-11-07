package fs

import (
	"ksef/core/ratelimits"
	"ksef/internal/utils"
	"os"
)

// read decodes the state file, tolerating a missing or empty file.
func (s *FSStorage) read() (ratelimits.State, error) {
	file, err := os.Open(s.dbFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer file.Close()

	var onDisk fsState
	if err = utils.ReadYAML(file, &onDisk); err != nil {
		return nil, err
	}

	if len(onDisk) == 0 {
		return nil, nil
	}

	state := make(ratelimits.State, len(onDisk))
	for operationId, slottedEntries := range onDisk {
		entries := make([]ratelimits.SlotEntries, len(slottedEntries))
		for idx, entry := range slottedEntries {
			entries[idx] = ratelimits.SlotEntries{
				Slot:    entry.Slot,
				Entries: entry.Entries,
			}
		}
		state[operationId] = entries
	}

	return state, nil
}
