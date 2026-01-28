package sessionregistry

import (
	"ksef/internal/utils"
	"path/filepath"
	"slices"
)

func (r *Registry) Save() error {
	if !r.dirty {
		return nil
	}

	defer func() { r.dirty = false }()

	// sort sessions by date in descending order so that the most recent
	// one ends up on the top.

	// theoretically, the session UID's are lexicographically sortable, but ..
	// nothing prevents the ministry of finance for this not to be the case
	slices.SortFunc(r.sessions, func(a, b *UploadSession) int {
		// if we'd use a.Timestamp.Compare then it would be in ascending order
		// which is not what we want.
		return b.Timestamp.Compare(a.Timestamp)
	})

	return utils.SaveYAML(r.sessions, filepath.Join(r.dir, registryFilename))
}
