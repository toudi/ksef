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

// ParseMonetaryValue takes a map and tries to parse the source data and return
// monetary value as a float that will be later given as input to MonetaryValue
// type.
func ParseMonetaryValue(data map[string]interface{}) (string, bool, error) {
	var parsedFloatNumber float64
	var is_a_float bool

	if value, exists := data["value"]; exists {
		// ok so this is potentially a serialized number. let's check if this is a valid one.
		tmpDecimalPlaces, exists := data["decimal-places"]
		if !exists {
			return "", false, fmt.Errorf("sub-struct does not contain decimal-places when value is a float")
		}
		decimalPlacesInt, ok := tmpDecimalPlaces.(int)
		if !ok {
			// decimal places is not an int.
			// could it be a string ?
			decimalPlacesString, ok := tmpDecimalPlaces.(string)
			var err error
			if ok {
				//
				decimalPlacesInt, err = strconv.Atoi(decimalPlacesString)
			} else {
				err = fmt.Errorf("sub-struct contains decimal-places but it's neither int nor string")
			}
			if err != nil {
				return "", false, fmt.Errorf("sub-struct does contain decimal-places but it is not an integer")
			}
		}
		decimalPlaces := decimalPlacesInt
		if parsedFloatNumber, is_a_float = value.(float64); is_a_float {
			return strconv.FormatFloat(parsedFloatNumber, 'f', decimalPlaces, 64), true, nil
		}
		if tmpInt, is_an_int := value.(int); is_an_int {
			parsedFloatNumber = float64(tmpInt) / math.Pow10(decimalPlaces)
		}

		if tmpString, is_a_string := value.(string); is_a_string {
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
			dotPosition := strings.Index(tmpString, ".")
			if decimalPlaces > 0 && dotPosition > 0 {
				stringLength := len(tmpString)
				if stringLength-dotPosition > decimalPlaces {
					tmpString = tmpString[:dotPosition+decimalPlaces+1]
				}
			}
			if dotPosition == -1 {
				tmpString = tmpString + ".0"
			}
			return tmpString, true, nil
		}

		return strconv.FormatFloat(parsedFloatNumber, 'f', decimalPlaces, 64), true, nil
	}
	// this is not a serialized decimal number, but it's also not a critical problem
	// could be just some other nested structure.
	return "", false, nil
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
	var err error

	value = ensureHasDot(value)

	unitPriceFloat, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return fmt.Errorf("cannot parse amount: %v", err)
	}

	m.DecimalPlaces = len(value) - strings.Index(value, ".") - 1
	m.Amount = int(unitPriceFloat * math.Pow10(m.DecimalPlaces))
	return nil
}
