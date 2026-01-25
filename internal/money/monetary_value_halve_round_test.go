package money_test

import (
	"ksef/internal/money"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHalveAndRound(t *testing.T) {
	type testCase struct {
		input    money.MonetaryValue
		expected money.MonetaryValue
	}

	for _, test := range []testCase{
		{
			// half of 0 is still 0
			input:    money.MonetaryValue{},
			expected: money.MonetaryValue{},
		},
		{
			input: money.MonetaryValue{
				Amount:        1005,
				DecimalPlaces: 2,
			},
			// 10.05 halved is 5.025, however after rounding up
			// it becomes 5.03
			expected: money.MonetaryValue{
				Amount:        503,
				DecimalPlaces: 2,
			},
		},
		{
			input: money.MonetaryValue{
				Amount:        1413,
				DecimalPlaces: 2,
			},
			// 14.13 halved is 7.065, however after rounding up
			// it becomes 7.07
			expected: money.MonetaryValue{
				Amount:        707,
				DecimalPlaces: 2,
			},
		},
		{
			input: money.MonetaryValue{
				Amount:        1020,
				DecimalPlaces: 2,
			},
			expected: money.MonetaryValue{
				Amount:        510,
				DecimalPlaces: 2,
			},
		},
		{
			input: money.MonetaryValue{
				Amount:        114,
				DecimalPlaces: 2,
			},
			expected: money.MonetaryValue{
				Amount:        57,
				DecimalPlaces: 2,
			},
		},
		{
			input: money.MonetaryValue{
				Amount:        1412,
				DecimalPlaces: 2,
			},
			// 14.12 halved is 7.06
			expected: money.MonetaryValue{
				Amount:        706,
				DecimalPlaces: 2,
			},
		},
		{
			input: money.MonetaryValue{
				Amount:        1406,
				DecimalPlaces: 2,
			},
			// 14.06 halved is 7.03
			expected: money.MonetaryValue{
				Amount:        703,
				DecimalPlaces: 2,
			},
		},
	} {
		t.Run(test.input.Format(2), func(t *testing.T) {
			t.Parallel()

			require.Equal(t, test.expected, test.input.HalveAndRoundUp())
		})
	}
}
