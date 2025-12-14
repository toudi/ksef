package invoicesdb

import (
	"bytes"
	"ksef/internal/certsdb"
	v2 "ksef/internal/client/v2"
	annualregistry "ksef/internal/invoicesdb/annual-registry"
	"ksef/internal/invoicesdb/config"
	monthlyregistry "ksef/internal/invoicesdb/monthly-registry"
	"ksef/internal/runtime"

	"github.com/spf13/viper"
)

// stateful invoices db registry
type InvoicesDB struct {
	cfg       config.InvoicesDBConfig
	importCfg config.ImportConfig

	contentBuffer   bytes.Buffer              // buffer for temporary XML
	monthlyRegistry *monthlyregistry.Registry // currently used monthly invoice registry
	annualRegistry  *annualregistry.Registry  // currently used annual invoice registry

	certsDB *certsdb.CertificatesDB // for retrieving offline certificate
	vip     *viper.Viper
	// for sync only
	prefix string
	// gateway client
	ksefClient *v2.APIClient
}

func NewInvoicesDB(vip *viper.Viper, initializers ...func(i *InvoicesDB)) (*InvoicesDB, error) {
	certsDB, err := certsdb.OpenOrCreate(runtime.GetGateway(vip))
	if err != nil {
		return nil, err
	}

	idb := &InvoicesDB{
		cfg:       config.GetInvoicesDBConfig(vip),
		importCfg: config.GetImportConfig(vip),
		certsDB:   certsDB,
		vip:       vip,
	}

	for _, initializer := range initializers {
		initializer(idb)
	}

	prefix, err := idb.getMonthlyRegistryPrefix()
	if err != nil {
		return nil, err
	}

	idb.prefix = prefix

	return idb, nil
}

func WithKSeFClient(client *v2.APIClient) func(idb *InvoicesDB) {
	return func(idb *InvoicesDB) {
		idb.ksefClient = client
	}
}

func (idb *InvoicesDB) Save() error {
	if idb.annualRegistry != nil {
		if err := idb.annualRegistry.Save(); err != nil {
			return err
		}
	}
	if idb.monthlyRegistry != nil {
		if err := idb.monthlyRegistry.Save(); err != nil {
			return err
		}
	}

	return nil
}
