package utils_test

import (
	"fmt"
	"ksef/internal/utils"
	"testing"
	"time"

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
