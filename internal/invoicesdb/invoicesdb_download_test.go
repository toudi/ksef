package invoicesdb

import (
	"fmt"
	"testing"
	"time"

	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
)

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

			// _ = generateMonthsRangeAtTime(test.today, test.startDate, test.endDate)

			require.Equal(
				t,
				test.expected,
				generateMonthsRangeAtTime(test.today, test.startDate, test.endDate),
			)
		})
	}
}
