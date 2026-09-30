package money

import "math"

// RoundedWhole returns a new MonetaryValue rounded to the nearest whole unit
// using the standard rounding rule (half rounded away from zero), with no
// decimal places. E.g. 123.3 -> 123, 123.5 -> 124, 390.67 -> 391.
// Used for JPK declaration fields (P_*) which must carry whole zloty amounts.
// Since the result has DecimalPlaces of 0, rendering it with FormatTrimmed
// produces a plain integer string.
func (m MonetaryValue) RoundedWhole() MonetaryValue {
	exponent := int(math.Pow10(m.DecimalPlaces))
	quotient := m.Amount / exponent
	remainder := m.Amount % exponent

	// Go's % takes the sign of the dividend, so negative amounts produce a
	// negative remainder. Take its magnitude to compare against the
	// rounding threshold.
	absRemainder := remainder
	if remainder < 0 {
		absRemainder = -remainder
	}

	// Standard rounding: half rounds away from zero. E.g. 123.50 (12350
	// groszy): quotient 123, remainder 50, and 50*2 >= 100 (the exponent),
	// so the quotient is bumped to 124; 123.49 has remainder 49 and
	// 49*2 < 100, so it stays 123. Negative amounts round away from zero
	// as well: -123.50 becomes -124.
	if absRemainder*2 >= exponent {
		if remainder < 0 {
			quotient--
		} else {
			quotient++
		}
	}

	return MonetaryValue{Amount: quotient}
}
