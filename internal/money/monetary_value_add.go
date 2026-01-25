package money

import (
	"math"
)

func (m MonetaryValue) Add(other MonetaryValue) MonetaryValue {
	if m.DecimalPlaces == other.DecimalPlaces {
		return MonetaryValue{
			Amount:        m.Amount + other.Amount,
			DecimalPlaces: m.DecimalPlaces,
		}
	}

	// if decimal places are different then we have to calculate a common and
	// pick the bigger one.

	// examples:
	// amount | decimal_places | actual number
	// 123.   | 2.             | 1.23.         (A)
	// 23.    | 0              | 23            (B)
	// 23     | 1              | 2.3           (C)
	//
	// A + B = 123 + 2300 = 2423 => 24.23
	// A + C = 123 + 230 = 353 => 3.53

	commonDecimalPlaces := m.DecimalPlaces
	if other.DecimalPlaces > commonDecimalPlaces {
		commonDecimalPlaces = other.DecimalPlaces
	}

	// fmt.Printf("common decimal places: %v\n", commonDecimalPlaces)
	// fmt.Printf("m.ToDecimalPlaces = %+v\n", m.ToDecimalPlaces(commonDecimalPlaces))
	// fmt.Printf("other.ToDecimalPlaces = %+v\n", other.ToDecimalPlaces(commonDecimalPlaces))

	return MonetaryValue{
		Amount:        m.ToDecimalPlaces(commonDecimalPlaces).Amount + other.ToDecimalPlaces(commonDecimalPlaces).Amount,
		DecimalPlaces: commonDecimalPlaces,
	}
}

func (m MonetaryValue) ToDecimalPlaces(places int) MonetaryValue {
	if m.DecimalPlaces == places {
		return m
	}

	var multiplier int = int(math.Pow10(places + m.DecimalPlaces))

	if m.DecimalPlaces < places {
		multiplier = int(math.Pow10(places - m.DecimalPlaces))
		return MonetaryValue{
			Amount:        m.Amount * multiplier,
			DecimalPlaces: places,
		}
	}

	return MonetaryValue{
		Amount:        m.Amount / multiplier,
		DecimalPlaces: places,
	}
}
