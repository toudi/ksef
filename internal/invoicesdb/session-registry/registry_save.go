package sessionregistry

import (
	"ksef/internal/utils"
	"path"
)

func (r *Registry) Save() error {
	if !r.dirty {
		return nil
	}

	defer func() { r.dirty = false }()

	return utils.SaveYAML(r.sessions, path.Join(r.dir, registryFilename))
}
