package invoicesdb

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"ksef/internal/client/v2/types/invoices"
	invoiceTypes "ksef/internal/client/v2/types/invoices"
	monthlyregistry "ksef/internal/invoicesdb/monthly-registry"
	"ksef/internal/logging"
	"ksef/internal/pdf"
	pdfconfig "ksef/internal/pdf/config"
	"ksef/internal/runtime"
	"ksef/internal/utils"
	"strings"
	"time"

	"github.com/spf13/viper"
)

var errDownloadingInvoices = errors.New("error downloading invoices")

func (i *InvoicesDB) DownloadInvoices(
	ctx context.Context,
	vip *viper.Viper,
	cfg invoiceTypes.DownloadParams,
) (err error) {
	// not sure if that's the "proper" way to do it, but let's just always use persistent time
	// and download invoices for all subjects. this way we can be incremental about it.
	tmpDownloadParams := invoiceTypes.DownloadParams{
		SubjectTypes:  cfg.SubjectTypes,
		PageSize:      cfg.PageSize,
		UseExportMode: cfg.UseExportMode,
	}

	// so that we know which registries to save
	affectedRegistries := make(map[*monthlyregistry.Registry]bool)
	lastTimestampPerRegistry := make(map[*monthlyregistry.Registry]time.Time)
	// just to be on the safe side - let's always try to download invoices for the
	// last month as well.
	monthsRange := generateMonthsRange(i.today, cfg.EndDate)

	if !cfg.StartDate.IsZero() {
		monthsRange = generateMonthsRange(cfg.StartDate, cfg.EndDate)
	}

	for idx, month := range monthsRange {
		registry, err := monthlyregistry.OpenOrCreateForMonth(
			vip,
			month,
		)
		if err != nil {
			return err
		}

		tmpDownloadParams.DateType = invoiceTypes.DateTypeStorage

		// initialize starting cutoff with month start:
		tmpDownloadParams.StartDate = month

		// if possible, override it with last sync timestamp:
		// (basically we're checking if the last timestamp is > start of month, however only if
		// it is within the month range. The reason for that is if for whatever reason there's
		// a wrong date persisted - i.e. from a different month due to a bug - we want to ignore
		// it as it would lead to a faulty request
		if registry.SyncParams != nil && utils.WithinMonthRange(
			registry.SyncParams.LastTimestamp, month,
		) {
			// use .In(month.Location) because the registry persists timestamps in UTC
			// however month comes from user entry, therefore we need to compare them with respect to
			// the same timezone
			// example: let's imagine that the registry saves something like
			// 2026-03-20T12:34:56 Z (UTC)
			// now, in order to get the first day of month this would be converted to
			// 2026-03-01T00:00:00 Z (UTC)
			// and now if user would pass just a date, without timezone, i.e. 2026-03-01
			// then WithinMonthRange could actually return false if the local timezone is positively offset
			// e.g. 2026-03-01T00:00:00 CET is actually 2026-02-28 UTC and therefore it would be out of range
			// if we'd compare it with the above calculated month start. Which is clearly not users intent.
			tmpDownloadParams.StartDate = registry.SyncParams.LastTimestamp.In(month.Location())
		}

		// last check. if the user explicitely requested a start date - let's prefer it (again, if it applies to the range)
		if utils.WithinMonthRange(cfg.StartDate, tmpDownloadParams.StartDate) && cfg.StartDate.Before(tmpDownloadParams.StartDate) {
			tmpDownloadParams.StartDate = cfg.StartDate
		}

		// get rid of end range initially ..
		tmpDownloadParams.EndDate = nil

		if idx < len(monthsRange)-1 {
			// unless there's a next month waiting to be processed - then we can easily determine
			// the end of the range.
			// ..
			// .. however because the API date range is inclusive - we've got to subtract 1 second
			// to end up at the last possible timestamp of `startDate` (ideally we should have actually
			// subtract less than that but the API clearly states that it accepts ISO-8601 which only
			// accepts full seconds)
			tmpTime := monthsRange[idx+1].Add(-time.Second)
			// because of the 1 second hack and because user can actually give us the last possible timestamp
			// let's make a sanity check over this - there's hardly any point for running the whole
			// download machinery for an empty range:
			if tmpTime.Equal(tmpDownloadParams.StartDate) {
				logging.DownloadLogger.Warn(
					"Pusty zakres dla pobierania faktur",
					"startDate", tmpDownloadParams.StartDate,
					"endDate", tmpTime,
				)
				continue
			}

			tmpDownloadParams.EndDate = &tmpTime
		}

		lastTimestampPerRegistry[registry] = registry.SyncParams.LastTimestamp

		downloader := i.ksefClient.InvoiceDownloader(
			vip,
			i.certsDB,
			tmpDownloadParams,
			registry,
		)

		if err = downloader.Download(
			ctx,
			func(subjectType invoices.SubjectType, invoice invoiceTypes.InvoiceMetadata, content bytes.Buffer) error {
				if registry.ContainsHash(invoice.Checksum()) {
					logging.DownloadLogger.Info("Ta faktura już znajduje się w rejestrze", "KSeFRefNo", invoice.KSeFNumber, "checksum", invoice.Checksum())
					return nil
				}
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
		if err = downloader.Close(); err != nil {
			return errors.Join(err, fmt.Errorf("błąd zamykania downloadera"))
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

func generateMonthsRange(startDate time.Time, endDate *time.Time) []time.Time {
	today := time.Now()

	return generateMonthsRangeAtTime(today, startDate, endDate)
}

func generateMonthsRangeAtTime(today time.Time, startDate time.Time, endDate *time.Time) []time.Time {
	// that is quite important - basically we do not control user input with regards to timezone.
	// theoretically we could - like we could just bail out with an error if they give an UTC
	// as this can get nasty with regards to last day + near midnight cases, but instead of
	// doing that let's just convert all of the dates to be within the same timezone.
	// if the user does not provide any timezone then .Now() will fallback to local timezone anyway
	today = today.In(startDate.Location())

	if endDate != nil && !endDate.IsZero() {
		today = (*endDate).In(startDate.Location())
	}

	monthsRange := []time.Time{}

	for !startDate.After(today) {
		monthsRange = append(monthsRange, startDate)

		// calculate the first day of next month.
		// In order to do that correctly we have to first override the day at `startDate` to 1. The
		// reason for that is - if somebody gives us a number that would be very close to the end
		// of the month (e.g. 25+) when we add +1 month we'd actually skip an entire month.
		// Therefore zeroing the time and overriding the day gives us safety
		startDate = time.Date(startDate.Year(), startDate.Month(), 1, 0, 0, 0, 0, startDate.Location())
		// add one month
		startDate = startDate.AddDate(0, 1, 0)
	}

	return monthsRange
}
