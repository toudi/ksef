package invoicesdb

import (
	"errors"
	"ksef/internal/invoice"
	annualregistry "ksef/internal/invoicesdb/annual-registry"
	monthlyregistry "ksef/internal/invoicesdb/monthly-registry"
	"ksef/internal/logging"
	"ksef/internal/sei"
	"ksef/internal/sei/generators/fa"
	"ksef/internal/utils"
	"time"
)

var errAddingItem = errors.New("unable to add item")

func (idb *InvoicesDB) handleCorrection(
	inv *sei.ParsedInvoice,
	originalInvoiceData *annualregistry.Invoice,
) (err error) {
	// create a timestamp of potential invoice generation. This will become important later
	generationTime := time.Now().Local()

	// in contrast to how KSeF calculates invoice checksum - here we can only calculate the checksum on the
	// actual content of what we're about to correct. The reason is very simple - because the import files
	// will be processed over and over again, we need to realize if the updated invoice (which we've already
	// discovered is a correction) is already in the system or does it need to be reported. But because
	// this code actually generates corrections (and that includes the correction number) and on top of that
	// relies on timestamps - we've got a chicken and the egg problem. However, if we hash only the contents
	// that has changed - we can calculate the checksum of that instead.
	correctionInvoice := &invoice.Invoice{
		Type: fa.InvoiceTypeCorrection,
		Correction: &invoice.CorrectionInfo{
			OriginalIssueDate: originalInvoiceData.GenerationTime,
			RefNo:             originalInvoiceData.RefNo,
			KSeFRefNo:         originalInvoiceData.KSeFRefNo,
		},
		Issuer:         inv.Invoice.Issuer,
		GenerationTime: generationTime,
		Attributes:     inv.Invoice.Attributes,
		Meta:           inv.Invoice.Meta,
		KSeFFlags:      &invoice.KSeFFlags{},
	}

	// we can only find the difference in items if the content is serialized
	if originalInvoiceData.Contents != "" {
		originalInvoice, err := originalInvoiceData.Unmarshall()
		if err != nil {
			return err
		}

		if len(inv.Invoice.Items) == 0 {
			// a special case - we're zeroing out the original invoice
			if err = correction_RemoveAllItems(
				originalInvoice,
				inv.Invoice,
				correctionInvoice,
			); err != nil {
				return err
			}
		} else if len(inv.Invoice.Items) > len(originalInvoice.Items) {
			// a new item has been added
			if err = correction_ItemHasBeenAdded(
				originalInvoice,
				inv.Invoice,
				correctionInvoice,
			); err != nil {
				return err
			}
		} else if len(inv.Invoice.Items) < len(originalInvoice.Items) {
			// item has been removed
			if err = correction_ItemHasBeenRemoved(
				originalInvoice,
				inv.Invoice,
				correctionInvoice,
			); err != nil {
				return err
			}
		} else {
			// number of items is equal
			if err = correction_NumberOfItemsEqual(
				originalInvoice,
				inv.Invoice,
				correctionInvoice,
			); err != nil {
				return err
			}
		}
	}

	for _, correction := range originalInvoiceData.Corrections {
		idb.contentBuffer.Reset()
		// temporarily replace generation time so that the content checksum can be deterministic
		correctionInvoice.GenerationTime = correction.GenerationTime
		correctionInvoice.Issued = correction.GenerationTime
		correctionInvoice.Number = correction.RefNo

		inv.Invoice = correctionInvoice
		if err := inv.ToXML(correction.GenerationTime, &idb.contentBuffer); err != nil {
			return err
		}

		correctionChecksum := utils.Sha256Hex(idb.contentBuffer.Bytes())

		if correction.Checksum == correctionChecksum {
			logging.GenerateLogger.Info("Korekta znaleziona w rejestrze faktur. No-op", "suma kontrolna", correctionChecksum)
			return nil
		}
	}

	idb.contentBuffer.Reset()

	// we were unable to find a matching correction. Let's replace the generation time for starters
	correctionInvoice.GenerationTime = generationTime
	correctionInvoice.Issued = generationTime

	// we can proceed with the correction generation.
	// let's generate correction number
	correctionNumber := idb.annualRegistry.GenerateCorrectionNumber(
		idb.subjectSettings.Import.AutoCorrection.Scheme,
		generationTime,
	)
	correctionInvoice.Number = correctionNumber

	// the reason why we're replacing the Invoice here is that in order to generate the XML
	// we also require issuer info, some form constants and so on which are all part of the
	// generator's commonData attribute. However, the generator is a protected variable so
	// it's just easiest to replace the invoice itself while leaving everything else as is.
	inv.Invoice = correctionInvoice

	if err := inv.ToXML(generationTime, &idb.contentBuffer); err != nil {
		return err
	}

	checksum := utils.Sha256Hex(idb.contentBuffer.Bytes())

	var monthlyRegistry *monthlyregistry.Registry

	if monthlyRegistry, err = idb.getMonthlyRegistryForInvoice(inv); err != nil {
		return err
	}

	fileName := monthlyRegistry.GetDestFileName(inv, monthlyregistry.InvoiceTypeIssued)
	if err = utils.SaveBufferToFile(idb.contentBuffer, fileName); err != nil {
		return errors.Join(errUnableToSaveInvoice, err)
	}

	if idb.importCfg.Offline || inv.Invoice.Issued.Before(idb.today) {
		inv.Invoice.KSeFFlags.Offline = true
	}

	if err := monthlyRegistry.AddInvoice(
		inv,
		monthlyregistry.InvoiceTypeIssued,
		checksum,
	); err != nil {
		return err
	}

	invoice := monthlyRegistry.GetInvoiceByChecksum(checksum)

	if err := originalInvoiceData.AddCorrection(
		correctionInvoice,
		checksum,
	); err != nil {
		return err
	}

	idb.newInvoices = append(idb.newInvoices, &NewInvoice{
		registry: monthlyRegistry,
		invoice:  invoice,
	})

	return nil
}
