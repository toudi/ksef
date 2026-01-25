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
