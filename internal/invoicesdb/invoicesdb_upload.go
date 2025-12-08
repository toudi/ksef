package invoicesdb

import (
	"context"
	"errors"
	"ksef/internal/invoicesdb/config"
	monthlyregistry "ksef/internal/invoicesdb/monthly-registry"
	"ksef/internal/invoicesdb/uploader"
	"ksef/internal/logging"
	"os"
	"path"
	"time"

	"github.com/spf13/viper"
)

func (i *InvoicesDB) UploadOutstandingInvoices(
	ctx context.Context,
	vip *viper.Viper,
) error {
	uploaderConfig := config.GetUploaderConfig(vip)

	var uploader = uploader.NewUploader(i.vip)
	// in order to handle the 31st day / 1st day problem, let's just
	// try to sync both of the months
	// (basically what this is all about is if you've imported invoices that were issued on the 31'st day
	// of the previous month but you're uploading them on the 1st day of the next month.

	var today = time.Now().Local()
	var previousMonth = today.AddDate(0, -1, 0)

	var months = []time.Time{previousMonth, today}

	for _, month := range months {
		// try to initialize monthly registry for the given month
		registry, err := monthlyregistry.OpenForMonth(i.prefix, month)
		if err != nil && os.IsNotExist(err) {
			logging.InvoicesDBLogger.Debug("registry does not exist; no-op", "dir", path.Join(i.prefix, month.Format("2006/01")))
			continue
		}
		if err != nil {
			return err
		}
		unsynced, err := registry.GetUnsynced()
		if err != nil {
			return err
		}
		for _, invoice := range unsynced {
			uploader.AddToQueue(invoice)
		}
	}

	if len(uploader.Queue) == 0 {
		logging.InvoicesDBLogger.Info("no unsynced invoices")
		return nil
	}

	_, uploadErr := uploader.UploadInvoices(ctx, uploaderConfig, i.ksefClient)
	if uploadErr != nil {
		return errors.Join(errUnableToUpload, uploadErr)
	}

	if uploaderConfig.DownloadUPO {
		// dispatch upo downloader
	}

	return nil
}
