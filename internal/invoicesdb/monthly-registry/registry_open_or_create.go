package monthlyregistry

import (
	"errors"
	"ksef/internal/certsdb"
	"ksef/internal/utils"
	"os"
	"path"

	"github.com/spf13/viper"
)

func OpenOrCreate(dir string, certsDB *certsdb.CertificatesDB, vip *viper.Viper) (*Registry, error) {
	regFile, exists, err := utils.FileExists(path.Join(dir, registryName))
	if err != nil && !os.IsNotExist(err) {
		// the only way for the err to be not nil is when there's a problem opening
		// file
		return nil, errOpeningRegistryFile
	}
	reg := &Registry{
		Invoices:   make([]*Invoice, 0),
		SyncParams: &SyncParams{},
		certsDB:    certsDB,
		vip:        vip,
		dir:        dir,
	}

	if exists {
		// if the file exists, then we need to read it's contents
		defer regFile.Close()
		if err = utils.ReadYAML(regFile, reg); err != nil {
			return nil, errors.Join(errReadingRegistryContents, err)
		}
	}

	return reg, nil
}
