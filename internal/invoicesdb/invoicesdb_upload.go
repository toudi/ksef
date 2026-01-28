package invoicesdb

import (
	"context"
	"errors"
	monthlyregistry "ksef/internal/invoicesdb/monthly-registry"
	statuschecker "ksef/internal/invoicesdb/status-checker"
	statuscheckerconfig "ksef/internal/invoicesdb/status-checker/config"
	"ksef/internal/invoicesdb/uploader"
	uploaderconfig "ksef/internal/invoicesdb/uploader/config"
	"ksef/internal/logging"
	"os"
	"path/filepath"
)

func (i *InvoicesDB) UploadOutstandingInvoices(
	ctx context.Context,
	uploaderConfig uploaderconfig.UploaderConfig,
	statusCheckerConfig statuscheckerconfig.StatusCheckerConfig,
) error {
	uploader := uploader.NewUploader(i.vip, uploaderConfig, i.ksefClient)
	// in order to handle the 31st day / 1st day problem, let's just
	// try to sync both of the months
	// (basically what this is all about is if you've imported invoices that were issued on the 31'st day
	// of the previous month but you're uploading them on the 1st day of the next month.

	months := i.monthsRange

	invoiceChecksumToRegistryMapping := make(map[string]*monthlyregistry.Registry)
	// these are all invoices that do not have KSeF number assigned right now, but
	// are being processed by some upload sessions.
	// we cannot upload them for the second time or it will result in an error, but we can
	// still check the status of unresolved sessions with the sync command
	invoiceChecksumsToSkip := make(map[string]bool)

	uploadSessionRegistry, err := i.getUploadSessionRegistry(i.today)
	if err != nil {
		return err
	}

	for _, month := range months {
		// try to initialize monthly registry for the given month
		registry, err := monthlyregistry.OpenForMonth(i.vip, month)
		if err != nil && os.IsNotExist(err) {
			logging.InvoicesDBLogger.Debug(
				"registry does not exist; no-op",
				"dir", filepath.Join(
					i.prefix, month.Format("2006"), month.Format("01"),
				),
			)
			continue
		}
		if err != nil {
			return err
		}
		unsynced, err := registry.GetUnsynced()
		if err != nil {
			return err
		}
		uploadSessionRegistryForMonth, err := i.getUploadSessionRegistry(month)
		if err != nil && !os.IsNotExist(err) {
			return err
		}
		// add any invoices that are pending to temporary map so that we don't accidentally
		// select them for upload for the second time.
		for _, hash := range uploadSessionRegistryForMonth.PendingUploadHashes() {
			invoiceChecksumsToSkip[hash] = true
		}
		for _, invoice := range unsynced {
			if invoiceChecksumsToSkip[invoice.Invoice.Checksum] {
				logging.InvoicesDBLogger.Info("faktura jest w trakcie przetwarzania - pomijam", "numer faktury", invoice.Invoice.RefNo, "suma kontrolna", invoice.Invoice.Checksum)
				continue
			}
			uploader.AddToQueue(invoice)
			invoiceChecksumToRegistryMapping[invoice.Invoice.Checksum] = registry
		}
	}

	if len(uploader.Queue) == 0 {
		logging.InvoicesDBLogger.Info("no unsynced invoices")
		return nil
	}

	// because we're forced to upload invoices in a group per form code, the uploader
	// returns a slice of upload sessions.
	uploadResult, uploadErr := uploader.UploadInvoices(ctx)
	if uploadErr != nil {
		return errors.Join(errUnableToUpload, uploadErr)
	}

	// with this pass of the Update function, we will persist information about upload session(s)
	// and write all assigned KSeF ref numbers (if any).
	// However, it's worth noting that these KSeF ref numbers do **NOT** represent the final
	// ref nos. Reason being that if there's something wrong with the invoice (e.g. it is an invalid
	// document or something like that) it won't be processed.
	// On top of that, if we're using a batch session then there won't be any ref no's
	// to begin with.
	for _, uploadSessionResult := range uploadResult {
		if err = uploadSessionRegistry.Update(
			uploadSessionResult,
			invoiceChecksumToRegistryMapping,
		); err != nil {
			return err
		}
	}

	if statusCheckerConfig.Wait {
		checker := statuschecker.NewStatusChecker(
			i.vip,
			i.ksefClient,
			statusCheckerConfig,
			i.monthsRange,
		)

		// this is slightly easier way since all of the upload sessions belong to the
		// same registry
		for _, uploadSessionResult := range uploadResult {
			checker.AddSessionID(uploadSessionResult.SessionID, uploadSessionRegistry)
		}

		checker.SetInvoiceHashToMonthlyRegistry(
			invoiceChecksumToRegistryMapping,
		)

		return checker.CheckSessions(ctx)
	}

	return nil
}
