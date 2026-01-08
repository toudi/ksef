package money

import (
	"errors"
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
func ParseMonetaryValue(data map[string]any) (string, bool, error) {
	var parsedFloatNumber float64
	var is_a_float bool

	if value, exists := data["value"]; exists {
		// ok so this is potentially a serialized number. let's check if this is a valid one.
		tmpDecimalPlaces, exists := data["decimal-places"]
		if !exists {
			return "", false, fmt.Errorf("sub-struct does not contain decimal-places when value is a float")
		}
		decimalPlaces, err := parseDecimal(tmpDecimalPlaces)
		if err != nil {
			return "", false, err
		}
		if parsedFloatNumber, is_a_float = value.(float64); is_a_float {
			return strconv.FormatFloat(parsedFloatNumber, 'f', decimalPlaces, 64), true, nil
		}
		tmpInt, err := parseDecimal(value)
		if err != nil {
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
			return "", false, err
		} else {
			parsedFloatNumber = float64(tmpInt) / math.Pow10(decimalPlaces)
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

func parseDecimal(input any) (result int, err error) {
	// let's start with uint64
	if tmpUint64, ok := input.(uint64); ok {
		// that went well..
		return int(tmpUint64), nil
	}

	// maybe it's an uint ?
	if tmpUint, ok := input.(uint); ok {
		return int(tmpUint), nil
	}

	if tmpInt, ok := input.(int); ok {
		return tmpInt, nil
	}

	// maybe it's a float64 if we're converting from JSON ?
	if tmpFloat64, ok := input.(float64); ok {
		return int(tmpFloat64), nil
	}

	return -1, errors.New("cannot recognize any of the types")
}
