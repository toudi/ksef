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
	// internal flags
	// prefix initialization is important during upload / sync commands
	// because we have to narrow down the invoice db to a single NIP passed
	// by the user. However, for import command we can rely on the NIP retrieved
	// from the invoice itself. Because we do not pass nip from the command line
	// and we do not actually use idb.prefix for anything during import - we can
	// safely skip this.
	skipPrefixInitialization bool
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

	if !idb.skipPrefixInitialization {
		prefix, err := idb.getMonthlyRegistryPrefix()
		if err != nil {
			return nil, err
		}

		idb.prefix = prefix
	}

	return idb, nil
}

func WithKSeFClient(client *v2.APIClient) func(idb *InvoicesDB) {
	return func(idb *InvoicesDB) {
		idb.ksefClient = client
	}
}

func WithoutInitializingPrefix() func(idb *InvoicesDB) {
	return func(idb *InvoicesDB) {
		idb.skipPrefixInitialization = true
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
