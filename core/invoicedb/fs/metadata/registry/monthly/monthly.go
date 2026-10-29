package monthly

import (
	"fmt"
	"ksef/core/invoicedb"
	"ksef/internal/utils"
	"path/filepath"
	"strconv"
)

const registryName = "registry.yaml"

type Registry struct {
	dbfile string

	state *state
}

func Open(root string, yearMonth invoicedb.YearMonth) (*Registry, error) {
	dbFileName := filepath.Join(root, strconv.Itoa(yearMonth.Year), fmt.Sprintf("%02d", yearMonth.Month), registryName)
	dbFile, _, err := utils.FileExists(dbFileName)
	if err != nil {
		return nil, err
	}
	defer dbFile.Close()

	var state state
	if err := utils.ReadYAML(dbFile, &state); err != nil {
		return nil, err
	}

	return &Registry{
		dbfile: dbFileName,
		state:  &state,
	}, nil
}

func (r *Registry) Close() error {
	return nil
}
