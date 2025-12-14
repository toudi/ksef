package invoicesdb

import (
	"errors"
	annualregistry "ksef/internal/invoicesdb/annual-registry"
	monthlyregistry "ksef/internal/invoicesdb/monthly-registry"
	sessionregistry "ksef/internal/invoicesdb/session-registry"
	"ksef/internal/runtime"
	"ksef/internal/sei"
	"os"
	"path"
	"time"
)

var (
	errUnableToCreateRegistryDir = errors.New("unable to create registry directory")
)

func (idb *InvoicesDB) getMonthlyRegistryPrefix() (string, error) {
	gateway := runtime.GetGateway(idb.vip)
	nip, err := runtime.GetNIP(idb.vip)
	if err != nil {
		return "", err
	}

	return path.Join(
		idb.cfg.Root,
		string(gateway),
		nip,
	), nil
}

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

func (idb *InvoicesDB) getUploadSessionRegistry(month time.Time) (*sessionregistry.Registry, error) {
	nip, err := runtime.GetNIP(idb.vip)
	if err != nil {
		return nil, err
	}

	gateway := runtime.GetGateway(idb.vip)

	// there is no active registry - let's try to create it.
	path := path.Join(
		idb.cfg.Root,
		string(gateway),
		nip,
		month.Format("2006"),
		month.Format("01"),
	)
	if err := os.MkdirAll(path, 0775); err != nil {
		return nil, errors.Join(errUnableToCreateRegistryDir, err)
	}

	return sessionregistry.OpenOrCreate(path)
}

func (idb *InvoicesDB) getUPODownloadPath(month time.Time) (string, error) {
	nip, err := runtime.GetNIP(idb.vip)
	if err != nil {
		return "", err
	}

	gateway := runtime.GetGateway(idb.vip)

	return path.Join(
		idb.cfg.Root,
		string(gateway),
		nip,
		month.Format("2006"),
		month.Format("01"),
		"upo",
	), nil
}
