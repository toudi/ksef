package common

import "math"

func RoundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

func AmountInGrosze(val float64) int {
	amountRounded := RoundFloat(val, 2)
	amountInGrosze := amountRounded * 100
	return int(amountInGrosze)
}
