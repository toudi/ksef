package monthlyregistry

// monthly registry represents a collection of invoices in a given month

import (
	"errors"
	"path/filepath"

	"ksef/internal/utils"
)

const (
	registryName = "registry.yaml"
)

var (
	errOpeningRegistryFile     = errors.New("error opening registry file")
	errReadingRegistryContents = errors.New("error reading registry contents")
)

func (r *Registry) Save() error {
	r.SavedOrdNums = r.OrdNums.ToSlice()
	return utils.SaveYAML(r, filepath.Join(r.dir, registryName))
}

func (r *Registry) Dir() string {
	return r.dir
}
