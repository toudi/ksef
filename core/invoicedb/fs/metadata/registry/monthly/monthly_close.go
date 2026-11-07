package monthly

import (
	"ksef/internal/utils"
	"os"
)

func (r *Registry) Close() error {
	if !r.dirty {
		return nil
	}

	r.state.SavedOrdNums = ordNumsToSlice(r.ordNum)

	if err := os.MkdirAll(r.dir, 0o755); err != nil {
		return err
	}
	if err := utils.SaveYAML(r.state, r.dbfile); err != nil {
		return err
	}

	r.dirty = false
	return nil
}
