package monthlyregistry

import (
	"errors"
	"ksef/internal/certsdb"
	"ksef/internal/invoicesdb/config"
	"ksef/internal/runtime"
	"ksef/internal/utils"
	"os"
	"path"
	"path/filepath"
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
		SyncParams:    &SyncParams{},
		checksumIndex: make(map[string]int),
		dir:           dir,
		vip:           vip,
		certsDB:       certsDB,
	}

	if exists {
		// if the file exists, then we need to read it's contents
		defer regFile.Close()
		var tmp Registry
		if err = utils.ReadYAML(regFile, &tmp); err != nil {
			return nil, errors.Join(errReadingRegistryContents, err)
		}

		// reading directly into &reg would cause any values that are not in the file
		// to be reverted to their zero values which in case of pointers means nil.
		// we definetely don't want that.

		if tmp.Invoices != nil {
			reg.Invoices = tmp.Invoices
		}
		if tmp.SyncParams != nil {
			reg.SyncParams = tmp.SyncParams
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

	registryPath := filepath.Join(
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
