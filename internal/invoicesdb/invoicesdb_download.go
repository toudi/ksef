package invoicesdb

import (
	"bytes"
	"context"
	"errors"
	"ksef/internal/client/v2/types/invoices"
	invoiceTypes "ksef/internal/client/v2/types/invoices"
	monthlyregistry "ksef/internal/invoicesdb/monthly-registry"
	"ksef/internal/pdf"
	pdfconfig "ksef/internal/pdf/config"
	"ksef/internal/runtime"
	"ksef/internal/utils"
	"strings"
	"time"

	"github.com/spf13/viper"
)

var errDownloadingInvoices = errors.New("error downloading invoices")

func (i *InvoicesDB) downloadInvoices(
	ctx context.Context,
	vip *viper.Viper,
	cfg invoiceTypes.DownloadParams,
) (err error) {
	// not sure if that's the "proper" way to do it, but let's just always use persistent time
	// and download invoices for all subjects. this way we can be incremental about it.
	tmpDownloadParams := invoiceTypes.DownloadParams{
		SubjectTypes: []invoiceTypes.SubjectType{
			invoiceTypes.SubjectTypeRecipient,
			invoiceTypes.SubjectTypePayer,
			invoiceTypes.SubjectTypeAuthorized,
		},
		PageSize:      cfg.PageSize,
		UseExportMode: cfg.UseExportMode,
	}

	// so that we know which registries to save
	affectedRegistries := make(map[*monthlyregistry.Registry]bool)
	lastTimestampPerRegistry := make(map[*monthlyregistry.Registry]time.Time)
	// just to be on the safe side - let's always try to download invoices for the
	// last month as well.
	monthsRange := i.monthsRange

	if !cfg.StartDate.IsZero() {
		monthsRange = generateMonthsRange(cfg.StartDate)
	}

	for idx, month := range monthsRange {
		registry, err := monthlyregistry.OpenOrCreateForMonth(
			vip,
			month,
		)
		if err != nil {
			return err
		}

		// initialize starting cutoff for downloading invoices
		tmpDownloadParams.StartDate = registry.SyncParams.LastTimestamp
		// get rid of end range initially ..
		tmpDownloadParams.EndDate = cfg.EndDate

		if idx < len(monthsRange)-2 {
			// unless there's a next month waiting to be processed - then we can easily determine
			// the end of the range.
			tmpDownloadParams.EndDate = &monthsRange[idx+1]
		}

		lastTimestampPerRegistry[registry] = registry.SyncParams.LastTimestamp

		downloader := i.ksefClient.InvoiceDownloader(
			tmpDownloadParams,
			registry,
		)

		if err = downloader.Download(
			ctx,
			func(subjectType invoices.SubjectType, invoice invoiceTypes.InvoiceMetadata, content bytes.Buffer) error {
				targetFilename := registry.GetDestFileNameForAPIInvoice(subjectType, invoice)
				if err = utils.SaveBufferToFile(content, targetFilename); err != nil {
					return err
				}
				if err = registry.AddReceivedInvoice(
					invoice,
					subjectType,
					runtime.GetEnvironment(i.vip),
				); err != nil {
					return err
				}

				if cfg.PDF {
					regInvoice := registry.GetInvoiceByChecksum(invoice.Checksum())
					printMeta := regInvoice.GetPrintingMeta()

					printer, err := pdf.GetInvoicePrinter(i.vip, pdfconfig.UsageSelector{
						Usage: printMeta.Usage,
					})
					if err != nil {
						return err
					}

					if err = printer.PrintInvoice(
						targetFilename,
						strings.Replace(targetFilename, ".xml", ".pdf", 1),
						printMeta,
					); err != nil {
						return err
					}
				}

				if invoice.StorageDate.After(lastTimestampPerRegistry[registry]) {
					lastTimestampPerRegistry[registry] = invoice.StorageDate
				}

				affectedRegistries[registry] = true

				return nil
			},
		); err != nil {
			return errors.Join(errDownloadingInvoices, err)
		}
	}

	for registry := range affectedRegistries {
		registry.SyncParams.LastTimestamp = lastTimestampPerRegistry[registry]

		if err = registry.Save(); err != nil {
			return err
		}
	}

	return nil
}

func generateMonthsRange(startDate time.Time) []time.Time {
	today := time.Now().Local()

	var monthsRange []time.Time
	for startDate.Before(today) {
		monthsRange = append(monthsRange, startDate)
		startDate = startDate.AddDate(0, 1, 0)
	}

	return monthsRange
}
