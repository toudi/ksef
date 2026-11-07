package monthly

import (
	"fmt"
	"ksef/core/invoicedb"
	"ksef/internal/utils"
	"os"
	"path/filepath"
	"strconv"
)

const registryName = "registry.yaml"

type Registry struct {
	dbfile    string
	dir       string
	yearMonth invoicedb.YearMonth

	state  *state
	ordNum map[InvoiceType]int

	checksumIndex map[string]int
	dirty         bool
}

func Open(root string, yearMonth invoicedb.YearMonth) (*Registry, error) {
	dir := filepath.Join(root, strconv.Itoa(yearMonth.Year), fmt.Sprintf("%02d", yearMonth.Month))
	dbFileName := filepath.Join(dir, registryName)

	registry := &Registry{
		dbfile:        dbFileName,
		dir:           dir,
		yearMonth:     yearMonth,
		state:         &state{},
		ordNum:        make(map[InvoiceType]int),
		checksumIndex: make(map[string]int),
	}

	dbFile, err := os.Open(dbFileName)
	if err != nil {
		if os.IsNotExist(err) {
			return registry, nil
		}
		return nil, err
	}
	defer dbFile.Close()

	if err := utils.ReadYAML(dbFile, registry.state); err != nil {
		return nil, err
	}

	savedOrdNums := len(registry.state.SavedOrdNums) > 0

	for index, invoice := range registry.state.Invoices {
		registry.checksumIndex[invoice.Checksum] = index
		if !savedOrdNums {
			registry.ordNum[invoice.Type]++
			invoice.OrdNum = registry.ordNum[invoice.Type]
		}
	}

	if savedOrdNums {
		registry.ordNum = registry.state.SavedOrdNums.ToMap()
	}

	return registry, nil
}
