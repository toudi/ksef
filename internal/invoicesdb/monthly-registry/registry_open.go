package monthlyregistry

import (
	"errors"
	"ksef/internal/invoicesdb/config"
	"ksef/internal/runtime"
	"ksef/internal/utils"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/viper"
)

func OpenForMonth(vip *viper.Viper, month time.Time) (*Registry, error) {
	var err error
	nip, err := runtime.GetNIP(vip)
	if err != nil {
		return nil, err
	}

	gateway := runtime.GetGateway(vip)
	invoicesDBConfig := config.GetInvoicesDBConfig(vip)

	registryPath := path.Join(
		invoicesDBConfig.Root,
		string(gateway),
		nip,
		month.Format("2006"),
		month.Format("01"),
	)
	regFilename := filepath.Join(
		registryPath,
		registryName,
	)

	return doOpen(regFilename)
}

func Open(path string) (*Registry, error) {
	regFilename := filepath.Join(path, registryName)

	return doOpen(regFilename)
}

func doOpen(regFilename string) (*Registry, error) {
	regFile, exists, err := utils.FileExists(regFilename)
	if err != nil && !os.IsNotExist(err) || !exists {
		return nil, err
	}

	reg := &Registry{
		Invoices:      make([]*Invoice, 0),
		SyncParams:    &SyncParams{},
		dir:           filepath.Dir(regFilename),
		OrdNums:       make(OrdNumsMap),
		SavedOrdNums:  make(OrdNums, 0),
		checksumIndex: make(map[string]int),
	}

	if exists {
		defer regFile.Close()

		if err = utils.ReadYAML(regFile, reg); err != nil {
			return nil, errors.Join(errReadingRegistryContents, err)
		}

		pathParts := strings.Split(reg.dir, string(filepath.Separator))
		pathLength := len(pathParts)

		if len(reg.SavedOrdNums) == 0 {
			reg.assignOrdNums()
		} else {
			reg.OrdNums = reg.SavedOrdNums.ToMap()
		}

		for index, invoice := range reg.Invoices {
			reg.checksumIndex[invoice.Checksum] = index
		}

		month, err := strconv.Atoi(pathParts[pathLength-1])
		if err != nil {
			return nil, err
		}
		year, err := strconv.Atoi(pathParts[pathLength-2])
		if err != nil {
			return nil, err
		}

		if reg.SyncParams.LastTimestamp.IsZero() {
			now := time.Now()
			reg.SyncParams.LastTimestamp = time.Date(
				year,
				time.Month(month),
				1,
				0, 0, 0, 0, now.Local().Location(),
			)
		}
	}

	return reg, nil
}
