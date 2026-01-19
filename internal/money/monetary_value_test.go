package money_test

import (
	"ksef/internal/money"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMonetaryValue(t *testing.T) {
	type testCase struct {
		input    string
		expected money.MonetaryValue
	}

	for _, test := range []testCase{
		{
			"",
			money.MonetaryValue{},
		},
		{
			"0",
			money.MonetaryValue{},
		},
		{
			".14",
			money.MonetaryValue{Amount: 14, DecimalPlaces: 2},
		},
		{
			"0.14",
			money.MonetaryValue{Amount: 14, DecimalPlaces: 2},
		},
		{
			"00000.14",
			money.MonetaryValue{Amount: 14, DecimalPlaces: 2},
		},
		{
			"0.140000",
			money.MonetaryValue{Amount: 14, DecimalPlaces: 2},
		},
		{
			"10",
			money.MonetaryValue{Amount: 10},
		},
		{
			"12.34",
			money.MonetaryValue{Amount: 1234, DecimalPlaces: 2},
		},
		{
			"12.",
			money.MonetaryValue{Amount: 12},
		},
	} {
		t.Run(test.input, func(t *testing.T) {
			var m money.MonetaryValue
			require.NoError(t, m.LoadFromString(test.input))
			require.Equal(t, test.expected, m)
		})
	}
}
