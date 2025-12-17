package invoicesdb

import (
	"context"
	"errors"
	"ksef/internal/client/v2/upo"
	"ksef/internal/http"
	"ksef/internal/invoicesdb/config"
	monthlyregistry "ksef/internal/invoicesdb/monthly-registry"
	"ksef/internal/invoicesdb/uploader"
	"ksef/internal/logging"
	"ksef/internal/pdf"
	"ksef/internal/pdf/printer"
	"os"
	"path"
	"strings"
	"time"

	"github.com/spf13/viper"
)

func (i *InvoicesDB) UploadOutstandingInvoices(
	ctx context.Context,
	vip *viper.Viper,
) error {
	uploaderConfig := config.GetUploaderConfig(vip)

	var uploader = uploader.NewUploader(i.vip, uploaderConfig, i.ksefClient)
	// in order to handle the 31st day / 1st day problem, let's just
	// try to sync both of the months
	// (basically what this is all about is if you've imported invoices that were issued on the 31'st day
	// of the previous month but you're uploading them on the 1st day of the next month.

	var today = time.Now().Local()
	var previousMonth = today.AddDate(0, -1, 0)

	var months = []time.Time{previousMonth, today}

	var invoiceChecksumToRegistryMapping = make(map[string]*monthlyregistry.Registry)

	uploadSessionRegistry, err := i.getUploadSessionRegistry(today)
	if err != nil {
		return err
	}

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
			invoiceChecksumToRegistryMapping[invoice.Invoice.Checksum] = registry
		}
	}

	if len(uploader.Queue) == 0 {
		logging.InvoicesDBLogger.Info("no unsynced invoices")
		return nil
	}

	uploadResult, uploadErr := uploader.UploadInvoices(ctx)
	if uploadErr != nil {
		return errors.Join(errUnableToUpload, uploadErr)
	}

	// with the interactive session, Update will actually already take care of
	// assigning KSeFRefNo to each of the invoices. However,
	// in batch mode, it will simply populate the timestamps and upload session
	// ref no's and nothing else
	if err = uploadSessionRegistry.Update(
		uploadResult,
		invoiceChecksumToRegistryMapping,
	); err != nil {
		return err
	}

	if uploaderConfig.WaitForStatus {
		if uploadResult, err = uploader.WaitForResult(ctx, uploadResult); err != nil {
			return err
		}

		// now the uploadResult will be updated with any potential errors and so on.
		if err = uploadSessionRegistry.Update(
			uploadResult, invoiceChecksumToRegistryMapping,
		); err != nil {
			return err
		}

		// we can only download UPO after we've waited for upload result
		if uploaderConfig.UPODownloader.Enabled {
			// dispatch upo downloader
			upoDestPath, err := i.getUPODownloadPath(today)
			if err != nil {
				return err
			}

			var printer printer.PDFPrinter

			if uploaderConfig.UPODownloader.ConvertToPDF {
				printer, err = pdf.GetLocalPrintingEngine()
				if err != nil {
					return err
				}
			}

			upoDownloader := upo.NewDownloader(
				http.NewClient(""), upo.UPODownloaderParams{
					Path:  upoDestPath,
					Mkdir: true,
					Wait:  uploaderConfig.UPODownloader.Timeout,
				},
			)

			for _, uploadSession := range uploadResult {
				if err = upoDownloader.Download(
					ctx,
					uploadSession.SessionID,
					uploadSession.Status.Upo.Pages,
					func(upoXMLFilename string) {
						if uploaderConfig.UPODownloader.ConvertToPDF {
							if err = printer.PrintUPO(
								upoXMLFilename, strings.Replace(upoXMLFilename, ".xml", ".pdf", 1),
							); err != nil {
								logging.PDFRendererLogger.Error("błąd konwersji UPO do PDF", "src", upoXMLFilename, "err", err)
							}
						}
					},
				); err != nil {
					return err
				}
			}
		}

	}

	return nil
}
