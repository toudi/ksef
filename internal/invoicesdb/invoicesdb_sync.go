package invoicesdb

import (
	"context"
	"errors"
	"ksef/internal/invoicesdb/config"
	"ksef/internal/logging"

	"github.com/spf13/viper"
)

var (
	errUnableToUpload           = errors.New("unable to upload invoices")
	errUnableToCheckSessions    = errors.New("unable to check sessions")
	errUnableToDownloadInvoices = errors.New("unable to download invoices")
	errUnableToGetSyncConfig    = errors.New("unable to get sync config")
)

func (i *InvoicesDB) Sync(ctx context.Context, vip *viper.Viper) error {
	var err error
	var syncConfig config.SyncConfig
	logger := logging.InvoicesDBLogger.With("tryb", "synchronizacja")

	if syncConfig, err = config.GetSyncConfig(vip); err != nil {
		return errors.Join(errUnableToGetSyncConfig, err)
	}

	origWaitForStatus := syncConfig.StatusCheckerConfig.Wait
	// forcefully turn off wait flag for the uploader part since the very
	// next thing is checking status for any pending sessions - which means that
	// the result of uploading invoices (if any) would also be taken care of
	syncConfig.StatusCheckerConfig.Wait = false

	// step 1. upload any invoices that are ready for processing
	logger.Info("wysyłam oczekujące faktury do KSeF")
	if err = i.UploadOutstandingInvoices(
		ctx,
		syncConfig.Uploader,
		syncConfig.StatusCheckerConfig,
	); err != nil {
		return errors.Join(errUnableToUpload, err)
	}

	syncConfig.StatusCheckerConfig.Wait = origWaitForStatus
	// step 2. check status for any upload sessions that have not yet been checked.
	logger.Info("sprawdzam stan oczekujących sesji")
	if err = i.checkPendingUploadSessions(
		ctx,
		syncConfig.StatusCheckerConfig,
	); err != nil {
		return errors.Join(errUnableToCheckSessions, err)
	}

	// step 3. download invoices
	logger.Info("pobieram faktury")
	if err = i.downloadInvoices(
		ctx,
		vip,
		syncConfig.Downloader,
	); err != nil {
		return errors.Join(errUnableToDownloadInvoices, err)
	}

	i.ksefClient.Close()

	return nil
}
