package invoicesdb

import (
	"context"
	"errors"
	"ksef/internal/certsdb"
	v2 "ksef/internal/client/v2"
	"ksef/internal/client/v2/auth/token"
	"ksef/internal/invoicesdb/config"
	monthlyregistry "ksef/internal/invoicesdb/monthly-registry"
	"ksef/internal/invoicesdb/uploader"
	"ksef/internal/logging"
	"ksef/internal/runtime"
	"os"
	"path"
	"time"

	"github.com/spf13/viper"
)

var (
	errUnableToUpload        = errors.New("unable to upload invoices")
	errUnableToGetSyncConfig = errors.New("unable to get sync config")
)

func (i *InvoicesDB) Sync(ctx context.Context, vip *viper.Viper) error {
	// in order to handle the 31st day / 1st day problem, let's just
	// try to sync both of the months
	// (basically what this is all about is if you've imported invoices that were issued on the 31'st day
	// of the previous month but you're uploading them on the 1st day of the next month.

	var today = time.Now().Local()
	var previousMonth = today.AddDate(0, -1, 0)

	var uploader = uploader.NewUploader(i.vip)

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

	gateway := runtime.GetGateway(vip)
	certsDB, err := certsdb.OpenOrCreate(gateway)
	if err != nil {
		return err
	}

	var clientInitializers = []v2.InitializerFunc{
		v2.WithAuthValidator(
			token.NewAuthHandler(
				vip,
				token.WithCertsDB(certsDB),
			),
		),
		v2.WithCertificatesDB(certsDB),
	}

	// we've got some invoices that need syncing. let's instantiate the client
	client, err := v2.NewClient(
		ctx,
		vip,
		clientInitializers...,
	)
	defer client.Close()

	syncConfig, err := config.GetSyncConfig(vip)
	if err != nil {
		return errors.Join(errUnableToGetSyncConfig, err)
	}
	uploaderConfig := syncConfig.Uploader

	_, uploadErr := uploader.UploadInvoices(ctx, uploaderConfig, client)
	if uploadErr != nil {
		return errors.Join(errUnableToUpload, err)
	}

	// great. all invoices have been pushed. now let's download invoices

	return nil
}
