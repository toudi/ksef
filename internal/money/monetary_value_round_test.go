package money_test

import (
	"ksef/internal/money"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRoundedWhole(t *testing.T) {
	type testCase struct {
		value    money.MonetaryValue
		expected money.MonetaryValue
	}

	for _, test := range []testCase{
		{value: money.MonetaryValue{}, expected: money.MonetaryValue{}},
		// round to nearest whole
		{value: money.MonetaryValue{Amount: 39067, DecimalPlaces: 2}, expected: money.MonetaryValue{Amount: 391}},
		{value: money.MonetaryValue{Amount: 446805, DecimalPlaces: 2}, expected: money.MonetaryValue{Amount: 4468}},
		// fractional below .5 rounds down, not up
		{value: money.MonetaryValue{Amount: 12330, DecimalPlaces: 2}, expected: money.MonetaryValue{Amount: 123}},
		{value: money.MonetaryValue{Amount: 12350, DecimalPlaces: 2}, expected: money.MonetaryValue{Amount: 124}},
		{value: money.MonetaryValue{Amount: 12349, DecimalPlaces: 2}, expected: money.MonetaryValue{Amount: 123}},
		{value: money.MonetaryValue{Amount: 12351, DecimalPlaces: 2}, expected: money.MonetaryValue{Amount: 124}},
		// exact integer
		{value: money.MonetaryValue{Amount: 1000, DecimalPlaces: 2}, expected: money.MonetaryValue{Amount: 10}},
		// negative amounts round half away from zero
		{value: money.MonetaryValue{Amount: -12330, DecimalPlaces: 2}, expected: money.MonetaryValue{Amount: -123}},
		{value: money.MonetaryValue{Amount: -12350, DecimalPlaces: 2}, expected: money.MonetaryValue{Amount: -124}},
		{value: money.MonetaryValue{Amount: -1167, DecimalPlaces: 2}, expected: money.MonetaryValue{Amount: -12}},
	} {
		t.Run(test.expected.FormatTrimmed(), func(t *testing.T) {
			t.Parallel()

			require.Equal(t, test.expected, test.value.RoundedWhole())
		})
	}
}
