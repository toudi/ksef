package invoicesdb

import (
	"bytes"
	"ksef/internal/certsdb"
	v2 "ksef/internal/client/v2"
	"ksef/internal/client/v2/session/status"
	annualregistry "ksef/internal/invoicesdb/annual-registry"
	"ksef/internal/invoicesdb/config"
	monthlyregistry "ksef/internal/invoicesdb/monthly-registry"
	"ksef/internal/runtime"
	"time"

	"github.com/spf13/viper"
)

type NewInvoice struct {
	registry *monthlyregistry.Registry
	invoice  *monthlyregistry.Invoice
}

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
	ksefClient    *v2.APIClient
	statusChecker *status.SessionStatusChecker
	// internal flags
	// prefix initialization is important during upload / sync commands
	// because we have to narrow down the invoice db to a single NIP passed
	// by the user. However, for import command we can rely on the NIP retrieved
	// from the invoice itself. Because we do not pass nip from the command line
	// and we do not actually use idb.prefix for anything during import - we can
	// safely skip this.
	skipPrefixInitialization bool
	// month range
	// we will be using this throughout couple of commands
	// basically all that it's about is to make sure that we cover the last day of previous
	// month / first day of current month scenario.
	// in other words, if some upload session / invoice was persisted on the last day
	// of the month but now we're past midnight, we want to take care of that as well
	monthsRange []time.Time
	today       time.Time
	// optimization for caching the filenames of offline invoices
	// for which we can generate PDF right away
	offlineInvoices []*NewInvoice
}

func newInvoicesDB(vip *viper.Viper) *InvoicesDB {
	// just so that we don't have to call time.Now() time and time again
	today := time.Now()
	previousMonth := today.AddDate(0, -1, 0)
	monthsRange := []time.Time{
		previousMonth, today,
	}

	idb := &InvoicesDB{
		cfg:         config.GetInvoicesDBConfig(vip),
		importCfg:   config.GetImportConfig(vip),
		vip:         vip,
		today:       today,
		monthsRange: monthsRange,
	}

	return idb
}

func NewInvoicesDB(vip *viper.Viper, initializers ...func(i *InvoicesDB)) (*InvoicesDB, error) {
	certsDB, err := certsdb.OpenOrCreate(runtime.GetGateway(vip))
	if err != nil {
		return nil, err
	}

	idb := newInvoicesDB(vip)
	idb.certsDB = certsDB

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
