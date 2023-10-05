package common

import (
	"math"
	"strconv"
)

func RoundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

func AmountInGrosze(val float64) int {
	amountRounded := RoundFloat(val, 2)
	amountInGrosze := amountRounded * 100
	return int(amountInGrosze)
}

func RenderFloatNumber(number float64) string {
	return strconv.FormatFloat(number, 'f', -1, 64)
}

func RenderAmountFromCurrencyUnits(amount int, decimalPlaces uint8) string {
	divisor := math.Pow(10, float64(decimalPlaces))
	return RenderFloatNumber(float64(amount) / divisor)
}
