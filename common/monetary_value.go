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
func (m *MonetaryValue) LoadFromString(value string, context map[string]string) error {
	// if the context contains number of decimal places then we can treat the number as integer
	var decimalPlaces string
	var exists bool
	var err error

	if decimalPlaces, exists = context["decimal-places"]; exists {
		if m.DecimalPlaces, err = strconv.Atoi(decimalPlaces); err != nil {
			return fmt.Errorf("%s does not seem to be an integer: %v", decimalPlaces, err)
		}
		// ok we have the number of decimal places so let's try to parse the number
		// itself.
		dotPosition := strings.Index(value, ".")
		// if there is no dot then the number is an integer encoded as string.
		if dotPosition == -1 {
			if m.Amount, err = strconv.Atoi(value); err != nil {
				return fmt.Errorf("unable to parse number from string %s: %v", value, err)
			}
			// everything is fine, the number was parsed without any problem.
			return nil
		}
		// if there is a dot then the number is a float encoded as string
		// let's check if there are more or equal number of decimal points in the string.
		// essentially let's consider this state table:
		//
		// raw value | decimal-places | expected value
		// ----------+----------------+---------------
		// 1.0       | 2              | 1.0
		// 1.23456   | 2              | 1.23
		//
		// so as you can see this is a way to force the parser to only read up until the number
		// of decimal places and make sure that the string was not created due to float rounding
		// error.
		stringLength := len(value)
		// if the string has less precision points than the expected one then we don't have to
		// do anything. Otherwise, let's truncate it.
		if stringLength-dotPosition > m.DecimalPlaces {
			value = value[:dotPosition+m.DecimalPlaces+1]
		}
	}
	// otherwise we have to take our best guess by treating the number as float.
	// this means that if it doesn't contain any decimal places at all we would
	// do
	//
	// 1 => 1.0
	//
	// and thus number of decimal places would be 1
	// if the number is already a float then we can parse number of decimal places
	// from there.
	value = ensureHasDot(value)

	unitPriceFloat, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return fmt.Errorf("cannot parse amount: %v", err)
	}

	m.DecimalPlaces = len(value) - strings.Index(value, ".") - 1
	m.Amount = int(unitPriceFloat * math.Pow10(m.DecimalPlaces))
	return nil
}
