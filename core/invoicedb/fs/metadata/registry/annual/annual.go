package annual

import (
	"ksef/core/invoicedb"
	"ksef/internal/logging"
	"ksef/internal/utils"
	"path/filepath"
	"strconv"
	"sync"
)

const registryName = "invoices.yaml"

var (
	openedRegistries map[string]map[int]*Registry
	mu               sync.Mutex
)

func init() {
	openedRegistries = make(map[string]map[int]*Registry)
}

type Registry struct {
	dbfile string

	invoices []*Invoice
}

func Open(root string, yearMonth invoicedb.YearMonth) (*Registry, error) {
	// root is only the prefix - i.e. data/{prod,demo,test}/{nip}
	mu.Lock()
	defer mu.Unlock()
	if registry, exists := openedRegistries[root][yearMonth.Year]; exists {
		logging.InvoicesDBLogger.Debug("found annual registry in cache - reusing", "root", root, "year", yearMonth.Year)
		return registry, nil
	}
	// if we're here then it means we've got to try to find the registry
	dbFileName := filepath.Join(root, strconv.Itoa(yearMonth.Year), registryName)
	dbFile, _, err := utils.FileExists(dbFileName)
	if err != nil {
		return nil, err
	}
	defer dbFile.Close()
	// TODO: parse db
	var invoices []*Invoice
	if err := utils.ReadYAML(dbFile, &invoices); err != nil {
		return nil, err
	}
	registry := &Registry{
		dbfile:   dbFileName,
		invoices: invoices,
	}

	if openedRegistries[root] == nil {
		openedRegistries[root] = make(map[int]*Registry)
	}

	openedRegistries[root][yearMonth.Year] = registry
	logging.InvoicesDBLogger.Debug("set annual registry to cache", "root", root, "year", yearMonth.Year)

	return registry, nil
}

func (r *Registry) Close() error {
	return nil
}
