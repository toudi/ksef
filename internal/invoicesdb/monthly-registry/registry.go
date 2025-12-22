package monthlyregistry

// monthly registry represents a collection of invoices in a given month

import (
	"errors"
	"ksef/internal/utils"
	"path"
)

const (
	registryName = "registry.yaml"
)

var (
	errOpeningRegistryFile     = errors.New("error opening registry file")
	errReadingRegistryContents = errors.New("error reading registry contents")
)

func (r *Registry) Save() error {
	return utils.SaveYAML(r, path.Join(r.dir, registryName))
}
