package client

import (
	"bytes"
	"context"
	"errors"
	"time"

	"ksef/core/downloader"
	"ksef/core/invoicedb"
	"ksef/core/ksef"
	"ksef/internal/client/v2/invoices"
	"ksef/internal/logging"
	"ksef/internal/utils"
)

var errDownloadingInvoices = errors.New("error downloading invoices")

func (c *Client) DownloadInvoices(
	ctx context.Context,
	cfg downloader.Config,
) error {
	monthsRange, err := c.monthsRange(ctx, cfg)
	if err != nil {
		return err
	}

	// Sync state is applied only after every month succeeds. On any error the
	// command fails, resources.Close skips InvoiceDB.Close, and nothing is
	// persisted.
	pendingSync := make(map[invoicedb.YearMonth]time.Time)

	dl := invoices.NewInvoiceDownloader(
		cfg,
		c.certsDB,
		c.authenticatedHTTPClient(),
		c.invoiceDB,
	)

	for idx, month := range monthsRange {
		yearMonth := invoicedb.YearMonth{
			Year:  month.Year(),
			Month: int(month.Month()),
		}

		syncParams, err := c.invoiceDB.GetSyncParams(ctx, yearMonth)
		if err != nil {
			return err
		}

		tmpDownloadParams := cfg
		tmpDownloadParams.DateType = ksef.DateTypeStorage
		// initialize starting cutoff with month start:
		tmpDownloadParams.StartDate = month

		// if possible, override it with last sync timestamp:
		// (basically we're checking if the last timestamp is > start of month, however only if
		// it is within the month range. The reason for that is if for whatever reason there's
		// a wrong date persisted - i.e. from a different month due to a bug - we want to ignore
		// it as it would lead to a faulty request
		if syncParams != nil && !syncParams.LastTimestamp.IsZero() && utils.WithinMonthRange(syncParams.LastTimestamp, month) {
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
			tmpDownloadParams.StartDate = syncParams.LastTimestamp.In(month.Location())
		}

		// last check. if the user explicitely requested a start date - let's prefer it (again, if it applies to the range)
		if !cfg.StartDate.IsZero() &&
			utils.WithinMonthRange(cfg.StartDate, tmpDownloadParams.StartDate) &&
			cfg.StartDate.Before(tmpDownloadParams.StartDate) {
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

		lastTimestamp := syncParams.LastTimestamp
		if lastTimestamp.IsZero() {
			// this could happen if the user has created some invoices first (i.e. Issued)
			// which would initialize the registry, but SyncParams would be empty
			lastTimestamp = month
		}

		err = dl.Download(
			ctx,
			tmpDownloadParams,
			func(subjectType ksef.SubjectType, meta ksef.InvoiceMetadata, content bytes.Buffer) error {
				exists, err := c.invoiceDB.ContainsHash(ctx, yearMonth, meta.Checksum())
				if err != nil {
					return err
				}
				if exists {
					logging.DownloadLogger.Info(
						"Ta faktura już znajduje się w rejestrze",
						"KSeFRefNo", meta.KSeFNumber,
						"checksum", meta.Checksum(),
					)
					return nil
				}

				stored, err := c.invoiceDB.StoreReceivedInvoice(
					ctx,
					yearMonth,
					subjectType,
					meta,
					bytes.NewReader(content.Bytes()),
					c.gateway.QRCodeBaseURL,
				)
				if err != nil {
					return err
				}

				if cfg.PDF {
					if c.pdfRenderer == nil {
						return errors.New("PDF rendering requested but no renderer is configured")
					}
					if err := c.pdfRenderer.RenderInvoice(ctx, stored); err != nil {
						return err
					}
				}

				if meta.StorageDate.After(lastTimestamp) {
					lastTimestamp = meta.StorageDate
					pendingSync[yearMonth] = lastTimestamp
				}
				return nil
			},
		)
		closeErr := dl.Close()
		if err != nil {
			return errors.Join(errDownloadingInvoices, err, closeErr)
		}
		if closeErr != nil {
			return errors.Join(errDownloadingInvoices, closeErr)
		}
	}

	for yearMonth, lastTimestamp := range pendingSync {
		if err := c.invoiceDB.SetSyncParams(
			ctx,
			yearMonth,
			invoicedb.InvoiceSyncParams{LastTimestamp: lastTimestamp},
		); err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) monthsRange(
	ctx context.Context,
	cfg downloader.Config,
) ([]time.Time, error) {
	// easiest case - the user has selected start date manually
	if !cfg.StartDate.IsZero() {
		return utils.GenerateMonthsRange(cfg.StartDate, cfg.EndDate), nil
	}

	now := time.Now().Local()
	var start time.Time

	// we need to determine start timestamp from the last sync state.
	syncState, err := c.invoiceDB.GetLastSyncState(ctx)
	if err != nil {
		return nil, err
	}

	if syncState != nil {
		start = time.Date(
			syncState.YearMonth.Year,
			time.Month(syncState.YearMonth.Month),
			1, 0, 0, 0, 0,
			now.Location(),
		)
	} else {
		// user is probably starting for the first time - assume the first day
		// of the current month.
		start = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		// just to be super safe - include the previous month as well.
		start = start.AddDate(0, -1, 0)
	}

	return utils.GenerateMonthsRange(start, cfg.EndDate), nil
}
