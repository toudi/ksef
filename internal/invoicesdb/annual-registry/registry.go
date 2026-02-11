package annualregistry

import (
	"errors"
	"ksef/internal/invoicesdb/config"
	"ksef/internal/runtime"
	"ksef/internal/utils"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
)

const (
	registryName = "invoices.yaml"
)

var (
	errOpeningRegistryFile     = errors.New("error opening registry file")
	errReadingRegistryContents = errors.New("error reading registry contents")
)

type Registry struct {
	invoices []*Invoice

	dir string
}

func OpenOrCreate(dir string) (*Registry, error) {
	regFile, exists, err := utils.FileExists(path.Join(dir, registryName))
	if err != nil && !os.IsNotExist(err) {
		// the only way for the err to be not nil is when there's a problem opening
		// file
		return nil, errors.Join(errOpeningRegistryFile, err)
	}

	reg := &Registry{
		invoices: make([]*Invoice, 0),
		dir:      dir,
	}

	if exists {
		// if the file exists, then we need to read it's contents
		defer regFile.Close()
		if err = utils.ReadYAML(regFile, &reg.invoices); err != nil {
			return nil, errors.Join(errReadingRegistryContents, err)
		}
	}

	return reg, nil
}

func (r *Registry) Save() error {
	return utils.SaveYAML(r.invoices, filepath.Join(r.dir, registryName))
}

func OpenForMonth(vip *viper.Viper, month time.Time) (*Registry, error) {
	var err error
	nip, err := runtime.GetNIP(vip)
	if err != nil {
		return nil, err
	}

	environmentId := runtime.GetEnvironmentId(vip)
	invoicesDBConfig := config.GetInvoicesDBConfig(vip)

	registryPath := filepath.Join(
		invoicesDBConfig.Root,
		environmentId,
		nip,
		month.Format("2006"),
	)

	return OpenOrCreate(registryPath)
}
