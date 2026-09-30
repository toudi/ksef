package money

import (
	"fmt"
	"math"
	"strconv"
)

func (m *MonetaryValue) Format(decimalPlaces int) string {
	adjusted := m.ToDecimalPlaces(decimalPlaces)

	exponent := int(math.Pow10(decimalPlaces))
	decimalPart := adjusted.Amount / exponent
	fractionalPart := adjusted.Amount % exponent

	if decimalPlaces > 0 {
		return fmt.Sprintf("%d.%0*d", decimalPart, decimalPlaces, fractionalPart)
	}
	return strconv.Itoa(decimalPart)
}

// FormatTrimmed renders the value using the minimum number of decimal places
// necessary (equivalent to Go's float 'f' format with precision -1). Trailing
// zeroes are dropped, so 93.00 renders as "93" and 93.73 as "93.73".
func (m MonetaryValue) FormatTrimmed() string {
	if m.DecimalPlaces == 0 {
		return strconv.Itoa(m.Amount)
	}

	return RenderAmountFromCurrencyUnits(m.Amount, uint8(m.DecimalPlaces))
}
