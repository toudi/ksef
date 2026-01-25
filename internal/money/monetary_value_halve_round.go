package money

import "math"

func (m MonetaryValue) HalveAndRoundUp() MonetaryValue {
	return MonetaryValue{
		Amount:        int(math.Ceil(float64(m.Amount) / 2)),
		DecimalPlaces: m.DecimalPlaces,
	}
}
