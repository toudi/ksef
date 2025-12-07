package invoicesdb

import (
	"errors"
	annualregistry "ksef/internal/invoicesdb/annual-registry"
	monthlyregistry "ksef/internal/invoicesdb/monthly-registry"
	"ksef/internal/runtime"
	"ksef/internal/sei"
	"os"
	"path"
)

var (
	errUnableToCreateRegistryDir = errors.New("unable to create registry directory")
)

func (idb *InvoicesDB) getMonthlyRegistryForInvoice(inv *sei.ParsedInvoice) (*monthlyregistry.Registry, error) {
	if idb.monthlyRegistry != nil {
		return idb.monthlyRegistry, nil
	}

	gateway := runtime.GetGateway(idb.vip)

	// there is no active registry - let's try to create it.
	path := path.Join(
		idb.cfg.Root,
		string(gateway),
		inv.Invoice.Issuer.NIP,
		inv.Invoice.Issued.Format("2006"),
		inv.Invoice.Issued.Format("01"),
	)
	if err := os.MkdirAll(path, 0775); err != nil {
		return nil, errors.Join(errUnableToCreateRegistryDir, err)
	}

	if reg, err := monthlyregistry.OpenOrCreate(path, idb.certsDB, idb.vip); err != nil {
		return nil, err
	} else {
		idb.monthlyRegistry = reg
	}

	return idb.monthlyRegistry, nil
}

func (idb *InvoicesDB) getAnnualRegistryForInvoice(inv *sei.ParsedInvoice) (*annualregistry.Registry, error) {
	if idb.annualRegistry != nil {
		return idb.annualRegistry, nil
	}

	gateway := runtime.GetGateway(idb.vip)

	// there is no active registry - let's try to create it.
	prefix := path.Join(
		idb.cfg.Root,
		string(gateway),
		inv.Invoice.Issuer.NIP,
		inv.Invoice.Issued.Format("2006"),
	)
	if err := os.MkdirAll(prefix, 0775); err != nil {
		return nil, errors.Join(errUnableToCreateRegistryDir, err)
	}

	if reg, err := annualregistry.OpenOrCreate(prefix); err != nil {
		return nil, err
	} else {
		idb.annualRegistry = reg
	}

	return idb.annualRegistry, nil
}
