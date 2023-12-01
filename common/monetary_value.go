package common

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func ensureHasDot(value string) string {
	if !strings.Contains(value, ".") {
		return value + ".0"
	}

	return value
}

type MonetaryValue struct {
	Amount        int
	DecimalPlaces int
}

// a brief explanation.
// raw value | normalized value | decimal index | string length | expected multiplier
// ----------+------------------+---------------+---------------+--------------------
// 1.001     | 1001             | 1             | 5             | 3
// 1.0       | 1                | 1             | 3             | 1
// 123.45    | 12345            | 3             | 6             | 2
// so as you can see the multiplier that would turn this raw float into an int is given
// as string length - index of decimal point - 1
func (m *MonetaryValue) LoadFromString(value string) error {
	value = ensureHasDot(value)

	unitPriceFloat, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return fmt.Errorf("cannot parse amount: %v", err)
	}

	m.DecimalPlaces = len(value) - strings.Index(value, ".") - 1
	m.Amount = int(unitPriceFloat * math.Pow10(m.DecimalPlaces))

	return nil
}
