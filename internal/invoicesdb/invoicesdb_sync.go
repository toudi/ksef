package invoicesdb

import (
	"context"
	"errors"

	"github.com/spf13/viper"
)

var (
	errUnableToUpload        = errors.New("unable to upload invoices")
	errUnableToGetSyncConfig = errors.New("unable to get sync config")
)

func (i *InvoicesDB) Sync(ctx context.Context, vip *viper.Viper) error {
	if err := i.UploadOutstandingInvoices(ctx, vip); err != nil {
		return errors.Join(errUnableToUpload, err)
	}

	// great. all invoices have been pushed. now let's download invoices

	return nil
}
