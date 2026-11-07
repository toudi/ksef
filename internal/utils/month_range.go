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
	monthStart = StartOfMonth(monthStart)
	// calculate the last possible time of the month
	// the -1 is so that we deduct 1 nanosecond and this will give us the last possible timestamp for the
	// calendar month of `monthStart`
	monthEnd := EndOfMonth(monthStart)

	if timestamp.Before(monthStart) || timestamp.After(monthEnd) {
		return false
	}

	return true
}

func StartOfMonth(timestamp time.Time) time.Time {
	return time.Date(timestamp.Year(), timestamp.Month(), 1, 0, 0, 0, 0, timestamp.Location())
}

func EndOfMonth(timestamp time.Time) time.Time {
	return timestamp.AddDate(0, 1, 0).Add(-1)
}

// GenerateMonthsRange returns the first day of every month from startDate
// through the current month (or endDate, if provided).
func GenerateMonthsRange(startDate time.Time, endDate *time.Time) []time.Time {
	return GenerateMonthsRangeAtTime(time.Now(), startDate, endDate)
}

// GenerateMonthsRangeAtTime is the testable variant of GenerateMonthsRange.
func GenerateMonthsRangeAtTime(today time.Time, startDate time.Time, endDate *time.Time) []time.Time {
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
