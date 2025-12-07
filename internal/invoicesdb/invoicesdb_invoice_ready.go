package invoicesdb

import (
	"errors"
	"fmt"
	annualregistry "ksef/internal/invoicesdb/annual-registry"
	"ksef/internal/logging"
	"ksef/internal/sei"
	"ksef/internal/utils"
	"time"
)

var (
	errUnableToRenderXML   = errors.New("unable to render XML to temporary buffer")
	errUnableToHash        = errors.New("unable to hash temporary invoice")
	ErrAutoCorrectDisabled = errors.New("auto-correct is disabled")
)

func (idb *InvoicesDB) InvoiceReady(inv *sei.ParsedInvoice) error {
	var err error
	var annualRegistry *annualregistry.Registry

	idb.contentBuffer.Reset()
	// the prefix will be in the following format:
	// data/1111111111/2025/
	//                      01/
	//                        registry.yaml
	//                        wystawione/
	//                        otrzymane/
	//                      02/
	//                      wysylki/01/
	//                      wysylki/01/2025-01-02T03:04:05/
	//                      invoices.yaml

	fmt.Printf("invoice: %+v\n", inv)
	if annualRegistry, err = idb.getAnnualRegistryForInvoice(inv); err != nil {
		return err
	}
	invoice := inv.Invoice

	// check if this invoice is already in the annual registry. if so, we can just no-op.
	// however because the checksum depends on generation time, we need to establish if
	// we have already rendered this invoice.
	// therefore let's go to the first step - check if invoice can be found by it's ref no.

	var lastKnownGenerationTimestamp time.Time
	// it will default to time.Zero, which means that when passed to ToXML() function
	// it won't get overriden, unless ..
	var _invoice *annualregistry.Invoice
	if _invoice = annualRegistry.GetByRefNo(invoice.Number); _invoice != nil {
		// .. unless we already have this invoice in our system.
		// the purpose of this trick is to always generate the same XML content so that
		// if we're parsing the same invoice (which has already been sent to KSeF) so that
		// we can detect this.
		lastKnownGenerationTimestamp = _invoice.GenerationTime
	}

	// now we can render the invoice to XML
	if err = inv.ToXML(lastKnownGenerationTimestamp, &idb.contentBuffer); err != nil {
		return errors.Join(errUnableToRenderXML, err)
	}

	// perfect. now let's calculate checksum
	checksum := utils.Sha256Hex(idb.contentBuffer.Bytes())

	// and now we can finally detect if this is potentially a correction candidate or
	// simply we've already processed this.
	if _invoice != nil && checksum == _invoice.Checksum {
		logging.GenerateLogger.Debug("faktura została już zaimportowana. no-op.")
		return nil
	}

	// is this a new invoice ? or..
	// the invoice wasn't sent yet. no problem - just override it
	if _invoice == nil || _invoice.KSeFRefNo == "" {
		return idb.handleNewInvoice(inv, annualRegistry, checksum)
	}

	// if not, it must be a correction
	if !idb.importCfg.AutoCorrection.Enabled {
		return ErrAutoCorrectDisabled
	}
	return idb.handleCorrection(inv, _invoice)
}
