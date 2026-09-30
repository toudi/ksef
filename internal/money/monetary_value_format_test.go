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

func TestFormatTrimmed(t *testing.T) {
	type testCase struct {
		value    money.MonetaryValue
		expected string
	}

	for _, test := range []testCase{
		{value: money.MonetaryValue{}, expected: "0"},
		{value: money.MonetaryValue{Amount: 93373, DecimalPlaces: 2}, expected: "933.73"},
		{value: money.MonetaryValue{Amount: 150, DecimalPlaces: 2}, expected: "1.5"},
		{value: money.MonetaryValue{Amount: 100, DecimalPlaces: 2}, expected: "1"},
		{value: money.MonetaryValue{Amount: 29255, DecimalPlaces: 2}, expected: "292.55"},
		{value: money.MonetaryValue{Amount: -1250, DecimalPlaces: 2}, expected: "-12.5"},
		{value: money.MonetaryValue{Amount: 500, DecimalPlaces: 1}, expected: "50"},
	} {
		t.Run(test.expected, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, test.expected, test.value.FormatTrimmed())
		})
	}
}
