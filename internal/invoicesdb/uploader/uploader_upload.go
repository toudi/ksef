package uploader

import (
	"context"
	"errors"
	v2 "ksef/internal/client/v2"
	monthlyregistry "ksef/internal/invoicesdb/monthly-registry"
	"ksef/internal/invoicesdb/uploader/config"
)

var (
	errUpload                         = errors.New("upload error")
	errCannotLookupSrcInvoice         = errors.New("this really should not happen. cannot lookup source invoice in the upload queue")
	errUpdatingInvoiceReferenceNumber = errors.New("unable to update uploaded invoice reference number")
	errSavingRegistry                 = errors.New("unable to save registry")
)

func (u *Uploader) UploadInvoices(ctx context.Context, config config.UploaderConfig, cli *v2.APIClient) (*UploadResult, error) {
	var err error
	var result *UploadResult

	var affectedRegistries = make(map[*monthlyregistry.Registry]bool)

	if config.WaitForStatus {
		if result, err = u.WaitForResult(ctx, config, cli); err != nil {
			return nil, err
		}

		// ok, now that waiting for status is complete, we can iterate through the upload result,
		// lookup affected invoices and update their reference numbers.
		for _, uploadedInvoices := range result.UploadSessions {
			for checksum, uploadStatus := range uploadedInvoices {
				// first let's lookup the corresponding monthly registry by the invoice's checksum
				registry, exists := u.registryByInvoiceChecksum[checksum]
				if !exists {
					return nil, errors.Join(errUpload, errCannotLookupSrcInvoice)
				}

				if err = registry.UpdateInvoiceByChecksum(checksum, func(invoice *monthlyregistry.Invoice) {
					invoice.KSeFRefNo = uploadStatus.KSeFRefNo
					invoice.UploadErrors = uploadStatus.Errors
				}); err != nil {
					return nil, errors.Join(errUpload, errUpdatingInvoiceReferenceNumber)
				}

				affectedRegistries[registry] = true
			}
		}

		// we need to persist the changes
		for registry := range affectedRegistries {
			if err = registry.Save(); err != nil {
				return nil, errors.Join(errUpload, errSavingRegistry)
			}
		}
	}
	return result, nil
}
