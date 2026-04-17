package utils

import (
	"time"
)

// WithinMonthRange checks that a `timestamp` is within range specified by
// `monthStart` and the last day of the same month. However, `timestamp` cannot be
// smaller than `monthStart`
func WithinMonthRange(timestamp time.Time, monthStart time.Time) bool {
	// let's reconstruct monthStart to be set to *actual* month start. The reason being - this function
	// will be used with timestamps loaded from invoices registries where the timestamps will fluctuate
	// and they won't always be set to actual month starts:
	monthStart = time.Date(monthStart.Year(), monthStart.Month(), 1, 0, 0, 0, 0, monthStart.Location())
	// calculate the last possible time of the month
	// the -1 is so that we deduct 1 nanosecond and this will give us the last possible timestamp for the
	// calendar month of `monthStart`
	lastMoment := monthStart.AddDate(0, 1, -1).Add(-1)

	if timestamp.Before(monthStart) || timestamp.After(lastMoment) {
		return false
	}

	return true
}
