package monthlyregistry

import (
	"errors"
	"ksef/internal/certsdb"
	"ksef/internal/invoicesdb/config"
	"ksef/internal/runtime"
	"ksef/internal/utils"
	"os"
	"path"
	"time"

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
		Invoices:      make([]*Invoice, 0),
		OrdNums:       make(OrdNumsMap),
		SavedOrdNums:  make(OrdNums, 0),
		SyncParams:    &SyncParams{},
		certsDB:       certsDB,
		vip:           vip,
		dir:           dir,
		checksumIndex: make(map[string]int),
	}

	if exists {
		// if the file exists, then we need to read it's contents
		defer regFile.Close()
		if err = utils.ReadYAML(regFile, reg); err != nil {
			return nil, errors.Join(errReadingRegistryContents, err)
		}

		if len(reg.SavedOrdNums) > 0 {
			reg.OrdNums = reg.SavedOrdNums.ToMap()
		}

		reg.postOpenHook()
	}

	return reg, nil
}

func OpenOrCreateForMonth(vip *viper.Viper, month time.Time) (*Registry, error) {
	nip, err := runtime.GetNIP(vip)
	if err != nil {
		return nil, err
	}

	environmentId := runtime.GetEnvironmentId(vip)
	invoicesDBConfig := config.GetInvoicesDBConfig(vip)

	registryPath := path.Join(
		invoicesDBConfig.Root,
		environmentId,
		nip,
		month.Format("2006"),
		month.Format("01"),
	)

	// note: this function does *NOT* set certsDB since it's only used for
	// downloading invoices
	registry, err := OpenOrCreate(registryPath, nil, vip)
	if err != nil {
		return nil, err
	}
	if registry.SyncParams.LastTimestamp.IsZero() {
		registry.SyncParams.LastTimestamp = month
	}
	return registry, nil
}
