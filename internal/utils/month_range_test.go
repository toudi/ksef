package utils_test

import (
	"fmt"
	"testing"
	"time"

	"ksef/internal/utils"

	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
)

func TestMonthRange(t *testing.T) {
	type testCase struct {
		timestamp  time.Time
		monthStart time.Time
		expected   bool
	}

	warsawTime, _ := time.LoadLocation("Europe/Warsaw")

	monthStart := time.Date(2026, 4, 1, 0, 0, 0, 0, warsawTime)
	monthStartUTC := time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC) // it is actually 2 hours later than `referenceTime`

	referenceTime := time.Date(2026, 4, 16, 14, 15, 16, 0, warsawTime)
	referenceTimeUTC := time.Date(2026, 4, 16, 14, 15, 16, 0, time.UTC)

	for _, test := range []testCase{
		{
			timestamp:  time.Time{},
			monthStart: monthStart,
			expected:   false,
		},
		{
			timestamp:  referenceTime,
			monthStart: monthStart,
			expected:   true,
		},
		{
			timestamp:  referenceTimeUTC,
			monthStart: monthStart,
			expected:   true,
		},
		{
			timestamp:  time.Date(2026, 5, 1, 0, 0, 0, 0, warsawTime),
			monthStart: monthStart,
			expected:   false,
		},
		{
			timestamp:  time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC),
			monthStart: monthStart,
			expected:   false,
		},
		{
			timestamp:  time.Date(2026, 3, 29, 0, 0, 0, 0, warsawTime),
			monthStart: monthStart,
			expected:   false,
		},
		{
			// this might seem counter-intuitive, however 23:59:59 UTC is actually the next day
			// in local timezone which is already within the month range
			timestamp:  time.Date(2026, 3, 31, 23, 59, 59, 59, time.UTC),
			monthStart: monthStart,
			expected:   true,
		},
		{
			// here's the next test case which is set in the same timezone and this one
			// returns false exactly as expected
			timestamp:  time.Date(2026, 3, 31, 23, 59, 59, 59, time.UTC),
			monthStart: monthStartUTC,
			expected:   false,
		},
		{
			timestamp:  monthStart,
			monthStart: referenceTime,
			expected:   true,
		},
		{
			timestamp:  referenceTime,
			monthStart: monthStartUTC,
			expected:   true,
		},
	} {
		t.Run(fmt.Sprintf("WithinMonthRange(%v, %v) == %v", test.timestamp.In(test.monthStart.Location()), test.monthStart, test.expected), func(t *testing.T) {
			t.Parallel()

			require.Equal(t, test.expected, utils.WithinMonthRange(test.timestamp, test.monthStart))
		})
	}
}

func TestStartOfMonth(t *testing.T) {
	type testCase struct {
		input    time.Time
		expected time.Time
	}

	warsawTime, _ := time.LoadLocation("Europe/Warsaw")

	for _, test := range []testCase{
		{
			input:    time.Time{},
			expected: time.Time{},
		},
		{
			input:    time.Date(2026, 4, 16, 14, 15, 16, 0, time.UTC),
			expected: time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			input:    time.Date(2026, 4, 16, 14, 15, 16, 0, warsawTime),
			expected: time.Date(2026, 4, 1, 0, 0, 0, 0, warsawTime),
		},
	} {
		t.Run(fmt.Sprintf("%v", test.input), func(t *testing.T) {
			t.Parallel()

			require.Equal(t, test.expected, utils.StartOfMonth(test.input))
		})
	}
}

func TestEndOfMonth(t *testing.T) {
	warsawTime, _ := time.LoadLocation("Europe/Warsaw")

	input := utils.StartOfMonth(time.Date(2026, 4, 29, 13, 14, 15, 16, warsawTime))
	lastMoment := utils.EndOfMonth(input)

	require.Equal(t, time.Date(2026, 4, 30, 23, 59, 59, 999999999, warsawTime), lastMoment)
}

func TestGenerateMonthsRange(t *testing.T) {
	// tests will always assume that `today` is set at 2026-04-01
	// it's supposed to be April's fools, however today I've discovered a bug at this function
	// so I guess the joke's on me.

	warsawTime, _ := time.LoadLocation("Europe/Warsaw")
	aprilsFools := time.Date(2026, 4, 1, 15, 14, 13, 12, warsawTime)

	type testCase struct {
		today     time.Time
		startDate time.Time
		endDate   *time.Time
		expected  []time.Time
	}

	for _, test := range []testCase{
		{
			today:     time.Date(2026, 4, 1, 0, 0, 0, 0, warsawTime),
			startDate: time.Date(2026, 3, 15, 0, 0, 0, 0, warsawTime),
			expected: []time.Time{
				time.Date(2026, 3, 15, 0, 0, 0, 0, warsawTime),
				time.Date(2026, 4, 1, 0, 0, 0, 0, warsawTime),
			},
		},
		{
			today:     time.Date(2026, 4, 1, 0, 0, 0, 0, warsawTime),
			startDate: time.Date(2026, 3, 15, 0, 0, 0, 0, warsawTime),
			endDate:   lo.ToPtr(time.Date(2026, 3, 31, 23, 59, 58, 57, warsawTime)),
			expected: []time.Time{
				time.Date(2026, 3, 15, 0, 0, 0, 0, warsawTime),
			},
		},
		{
			today:     aprilsFools,
			startDate: time.Date(2026, 3, 15, 0, 0, 0, 0, warsawTime),
			expected: []time.Time{
				time.Date(2026, 3, 15, 0, 0, 0, 0, warsawTime),
				time.Date(2026, 4, 1, 0, 0, 0, 0, warsawTime),
			},
		},
		{
			today:     aprilsFools,
			startDate: time.Date(2026, 2, 15, 0, 0, 0, 0, warsawTime),
			expected: []time.Time{
				time.Date(2026, 2, 15, 0, 0, 0, 0, warsawTime),
				time.Date(2026, 3, 1, 0, 0, 0, 0, warsawTime),
				time.Date(2026, 4, 1, 0, 0, 0, 0, warsawTime),
			},
		},
		{
			today:     time.Date(2026, 4, 13, 12, 11, 10, 9, warsawTime),
			startDate: time.Date(2026, 2, 15, 0, 0, 0, 0, warsawTime),
			expected: []time.Time{
				time.Date(2026, 2, 15, 0, 0, 0, 0, warsawTime),
				time.Date(2026, 3, 1, 0, 0, 0, 0, warsawTime),
				time.Date(2026, 4, 1, 0, 0, 0, 0, warsawTime),
			},
		},
		{
			today:     time.Date(2026, 4, 11, 12, 11, 10, 9, warsawTime),
			startDate: time.Date(2026, 1, 31, 23, 59, 59, 0, time.UTC),
			expected: []time.Time{
				time.Date(2026, 1, 31, 23, 59, 59, 0, time.UTC),
				time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC),
			},
		},
	} {
		t.Run(fmt.Sprintf("generateMonthsRange(%v, %v)", test.startDate, test.endDate), func(t *testing.T) {
			t.Parallel()

			// _ = utils.GenerateMonthsRangeAtTime(test.today, test.startDate, test.endDate)

			require.Equal(
				t,
				test.expected,
				utils.GenerateMonthsRangeAtTime(test.today, test.startDate, test.endDate),
			)
		})
	}
}
