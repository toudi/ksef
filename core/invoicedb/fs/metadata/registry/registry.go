package registry

import (
	"errors"
	"ksef/core/invoicedb"
	"ksef/core/invoicedb/fs/metadata/registry/annual"
	"ksef/core/invoicedb/fs/metadata/registry/monthly"
	"ksef/internal/logging"
)

type FSInvoiceRegistry struct {
	root      string
	yearMonth invoicedb.YearMonth

	annualRegistry  *annual.Registry
	monthlyRegistry *monthly.Registry
}

func Open(root string, yearMonth invoicedb.YearMonth) (*FSInvoiceRegistry, error) {
	logging.InvoicesDBLogger.Debug("fs/metadata/registry::Open", "root", root, "yearMonth", yearMonth)
	annualRegistry, err := annual.Open(root, yearMonth)
	if err != nil {
		return nil, err
	}

	monthlyRegistry, err := monthly.Open(root, yearMonth)
	if err != nil {
		return nil, err
	}

	return &FSInvoiceRegistry{
		root:            root,
		yearMonth:       yearMonth,
		annualRegistry:  annualRegistry,
		monthlyRegistry: monthlyRegistry,
	}, nil
}

func (fs *FSInvoiceRegistry) Close() error {
	return errors.Join(fs.monthlyRegistry.Close(), fs.annualRegistry.Close())
}
