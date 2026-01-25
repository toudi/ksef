package money

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMonetaryValueAdd(t *testing.T) {
	type testCase struct {
		value    MonetaryValue
		other    MonetaryValue
		expected MonetaryValue
	}

	for _, test := range []testCase{
		{
			value:    MonetaryValue{Amount: 123, DecimalPlaces: 2},
			other:    MonetaryValue{Amount: 23},
			expected: MonetaryValue{Amount: 2423, DecimalPlaces: 2},
		},
		{
			value:    MonetaryValue{Amount: 123, DecimalPlaces: 2},
			other:    MonetaryValue{Amount: 23, DecimalPlaces: 1},
			expected: MonetaryValue{Amount: 353, DecimalPlaces: 2},
		},
		{
			value:    MonetaryValue{Amount: 123, DecimalPlaces: 1},
			other:    MonetaryValue{Amount: 23, DecimalPlaces: 2},
			expected: MonetaryValue{Amount: 1253, DecimalPlaces: 2},
		},
	} {
		require.Equal(t, test.expected, test.value.Add(test.other))
	}
}
