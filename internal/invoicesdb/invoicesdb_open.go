package invoicesdb

import (
	"errors"
	"ksef/internal/invoicesdb/config"
	"ksef/internal/runtime"
	"os"
	"path"

	"github.com/spf13/viper"
)

var (
	errInvalidRegistry = errors.New("specified invoices registry does not exist")
)

func OpenForNIP(nip string, vip *viper.Viper, initializers ...func(*InvoicesDB)) (*InvoicesDB, error) {
	cfg := config.GetInvoicesDBConfig(vip)
	gateway := runtime.GetGateway(vip)

	// this prefix does not contain months yet - it is the entrypoint for further processing
	// (like uploading invoices)
	var prefix = path.Join(
		cfg.Root,
		string(gateway),
		nip,
	)

	if _, err := os.Stat(prefix); err != nil && os.IsNotExist(err) {
		return nil, errors.Join(errInvalidRegistry, err)
	}

	idb := newInvoicesDB(vip)
	idb.prefix = prefix

	for _, initializer := range initializers {
		initializer(idb)
	}

	return idb, nil
}
