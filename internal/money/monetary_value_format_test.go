package money_test

import (
	"ksef/internal/money"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRenderingMonetaryValues(t *testing.T) {
	type testCase struct {
		value         money.MonetaryValue
		decimalPlaces int
		expected      string
	}

	for _, test := range []testCase{
		{
			value:    money.MonetaryValue{},
			expected: "0",
		},
		{
			value:         money.MonetaryValue{Amount: 1234, DecimalPlaces: 2},
			decimalPlaces: 2,
			expected:      "12.34",
		},
		{
			value:         money.MonetaryValue{Amount: 1234, DecimalPlaces: 2},
			decimalPlaces: 4,
			expected:      "12.3400",
		},
		{
			value:         money.MonetaryValue{Amount: 1200, DecimalPlaces: 2},
			decimalPlaces: 2,
			expected:      "12.00",
		},
		{
			value:         money.MonetaryValue{Amount: -1200, DecimalPlaces: 2},
			decimalPlaces: 2,
			expected:      "-12.00",
		},
	} {
		t.Run(test.expected, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, test.expected, test.value.Format(test.decimalPlaces))
		})
	}
}
