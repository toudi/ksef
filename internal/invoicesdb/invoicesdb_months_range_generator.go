package invoicesdb

import (
	monthlyregistry "ksef/internal/invoicesdb/monthly-registry"
	"ksef/internal/utils"
	"time"

	"github.com/spf13/viper"
)

// MonthsRangeGenerator is a function type that generates the months range
// for invoice operations.
type MonthsRangeGenerator func(vip *viper.Viper, today time.Time) []time.Time

// monthsRangeLastTimestampGenerator checks the last 12 months of registries
// and returns a months range starting from the most recent month with a
// non-zero LastTimestamp. Falls back to previous month if no such registry
// is found.
func MonthsRangeLastKnownTimestampGenerator(vip *viper.Viper, today time.Time) []time.Time {
	// Check the last 12 months, starting from the current month going backwards
	month := time.Date(today.Year(), today.Month(), 1, 0, 0, 0, 0, today.Location())

	for range 12 {
		// Try to open the registry for this month
		reg, err := monthlyregistry.OpenForMonth(vip, month)
		// a bit ugly hack / antipattern with checking for err == nil, however we do want to
		// make sure that if the for-loop iterates over we'll already have previous month.
		// otherwise we'd be perpetually checking the same month
		if reg != nil && err == nil {
			// If LastTimestamp is non-zero, use this month as the start
			if !reg.SyncParams.LastTimestamp.IsZero() {
				return generateMonthsRange(month, &today)
			}
		}

		month = month.AddDate(0, -1, 0)
	}

	// we were not able to establish the last known timestamp. fallback to previous month instead.

	// Fallback: use previous month
	return MonthsRangePreviousMonthGenerator(vip, today)
}

// monthsRangePreviousMonthGenerator generates a months range for the previous
// month only. This is suitable for operations like uploading invoices.
func MonthsRangePreviousMonthGenerator(vip *viper.Viper, today time.Time) []time.Time {
	previousMonth := utils.StartOfMonth(today).AddDate(0, -1, 0)
	return generateMonthsRange(previousMonth, &today)
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
