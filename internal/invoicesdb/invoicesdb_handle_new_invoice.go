package invoicesdb

import (
	"errors"
	annualregistry "ksef/internal/invoicesdb/annual-registry"
	monthlyregistry "ksef/internal/invoicesdb/monthly-registry"
	"ksef/internal/sei"
	"ksef/internal/utils"
)

var (
	errUnableToSaveInvoice = errors.New("unable to save invoice to a file")
)

func (idb *InvoicesDB) handleNewInvoice(
	inv *sei.ParsedInvoice,
	annualRegistry *annualregistry.Registry,
	checksum string,
) error {
	var monthlyRegistry *monthlyregistry.Registry
	var err error

	if monthlyRegistry, err = idb.getMonthlyRegistryForInvoice(inv); err != nil {
		return err
	}

	fileName := monthlyRegistry.GetDestFileName(inv, monthlyregistry.InvoiceTypeIssued)
	if err = utils.SaveBufferToFile(idb.contentBuffer, fileName); err != nil {
		return errors.Join(errUnableToSaveInvoice, err)
	}

	if idb.importCfg.Offline {
		inv.Invoice.KSeFFlags.Offline = true
	}
	// let's add information about the invoice to the monthly registry
	if err = monthlyRegistry.AddInvoice(
		inv,
		monthlyregistry.InvoiceTypeIssued,
		checksum,
	); err != nil {
		return err
	}

	// now let's save info about the invoice to the annual registry
	return annualRegistry.AddInvoice(
		inv,
		checksum,
		idb.importCfg.AutoCorrection.Enabled,
	)
}
