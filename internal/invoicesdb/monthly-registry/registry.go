package monthlyregistry

// monthly registry represents a collection of invoices in a given month

import (
	"errors"
	"ksef/internal/certsdb"
	"ksef/internal/utils"
	"os"
	"path"
	"time"

	"github.com/spf13/viper"
)

const (
	registryName = "registry.yaml"
)

var (
	errOpeningRegistryFile     = errors.New("error opening registry file")
	errReadingRegistryContents = errors.New("error reading registry contents")
)

func OpenOrCreate(dir string, certsDB *certsdb.CertificatesDB, vip *viper.Viper) (*Registry, error) {
	regFile, exists, err := utils.FileExists(path.Join(dir, registryName))
	if err != nil && !os.IsNotExist(err) {
		// the only way for the err to be not nil is when there's a problem opening
		// file
		return nil, errOpeningRegistryFile
	}
	var reg = &Registry{
		invoices:       make([]*Invoice, 0),
		uploadSessions: make([]*UploadSession, 0),
		certsDB:        certsDB,
		vip:            vip,
		dir:            dir,
	}

	if exists {
		// if the file exists, then we need to read it's contents
		defer regFile.Close()
		if err = utils.ReadYAML(regFile, &reg); err != nil {
			return nil, errors.Join(errReadingRegistryContents, err)
		}
	}

	return reg, nil
}

func OpenForMonth(prefix string, month time.Time) (*Registry, error) {
	var err error

	var registryPath = path.Join(
		prefix,
		month.Format("2006"),
		month.Format("01"),
	)
	var regFilename = path.Join(
		registryPath,
		registryName,
	)

	regFile, exists, err := utils.FileExists(regFilename)
	if err != nil && !os.IsNotExist(err) || !exists {
		return nil, err
	}

	var reg = &Registry{
		invoices: make([]*Invoice, 0),
		dir:      registryPath,
	}

	if exists {
		defer regFile.Close()

		if err = utils.ReadYAML(regFile, &reg.invoices); err != nil {
			return nil, errors.Join(errReadingRegistryContents, err)
		}
	}

	return reg, nil
}

func (r *Registry) Save() error {
	return utils.SaveYAML(r.invoices, path.Join(r.dir, registryName))
}
