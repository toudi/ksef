package inputprocessors

import (
	"fmt"
	"ksef/common"
	"strconv"
	"strings"
)

func (y *YAMLDecoder) processSingleInvoiceSource(invoice map[string]interface{}, parser *common.Parser) error {
	if err := processRecurse("", invoice, parser); err != nil {
		return fmt.Errorf("error running processRecurse: %v", err)
	}
	return parser.InvoiceReady()
}

func processRecurse(section string, data interface{}, parser *common.Parser) error {
	var err error
	var fullSection string
	var hasScalarData bool = false
	var sectionData map[string]interface{}
	// because # is used to denote namespace attributes it would mean that a yaml key
	// would start with a hash, but that's not a valid YAML. therefore instead of
	// writing
	//
	// #foo: bar
	//
	// we use a notation of
	//
	// \#foo: bar
	//
	// similarly, instead of writing
	//
	// #foo:bar: baz
	//
	// which is technically still a valid yaml but may be confusing, the following
	// syntax is used:
	//
	// \#foo\:bar: baz
	//
	var keyNameReplacer *strings.Replacer = strings.NewReplacer("\\#", "#", "\\:", ":")

	if keyData, is_a_map := data.(map[string]interface{}); is_a_map {
		sectionData = make(map[string]interface{})

		for keyName, keyDataItem := range keyData {
			fullSection = keyName
			if section != "" {
				fullSection = section + "." + keyName
			}
			// let's check if this data is also an interface
			if _tmpData, is_a_map := keyDataItem.(map[string]interface{}); is_a_map {
				parsed_number, is_serialized_number, err := common.ParseMonetaryValue(_tmpData)
				if err != nil {
					return fmt.Errorf("error during parseSerializedNumber(): %v", err)
				}
				if is_serialized_number {
					sectionData[keyName] = parsed_number
					hasScalarData = true
					continue
				}
				if err = processRecurse(fullSection, keyDataItem, parser); err != nil {
					return fmt.Errorf("error in recursing to processRecurse(): %v", err)
				}
			} else if _, is_an_array := keyDataItem.([]interface{}); is_an_array {
				for _, arrayElement := range keyDataItem.([]interface{}) {
					if err = processRecurse(fullSection, arrayElement, parser); err != nil {
						return fmt.Errorf("error in recursing to processRecurse() for an array item: %v", err)
					}
				}
			} else {
				// this is a scalar therefore we can safely store the values.
				sectionData[keyNameReplacer.Replace(keyName)] = keyDataItem
				hasScalarData = true
			}
		}

	}

	if hasScalarData {
		// fmt.Printf("scalar data for section %s:\n%+v\n", section, sectionData)

		// let's yield 3 lines. first would be the section name:
		if err = parser.ProcessLine([]string{"Sekcja", section}); err != nil {
			return fmt.Errorf("error processing headers: %v", err)
		}

		// next, the section headers
		sectionDataHeaders := make([]string, len(sectionData))
		// because of random order of iteration over the map values, let's just create
		// second array and use a single forloop to populate both of them
		sectionDataValues := make([]string, len(sectionData))

		var fieldIndex int = 0
		for keyName, scalarValue := range sectionData {
			sectionDataHeaders[fieldIndex] = keyName

			if scalarValueString, is_a_string := scalarValue.(string); is_a_string {
				sectionDataValues[fieldIndex] = scalarValueString
			} else if scalarValueInt, is_an_int := scalarValue.(int); is_an_int {
				sectionDataValues[fieldIndex] = strconv.Itoa(scalarValueInt)
			} else if scalarValueFloat, is_a_float := scalarValue.(float64); is_a_float {
				sectionDataValues[fieldIndex] = strconv.FormatFloat(scalarValueFloat, 'f', -1, 64)
			} else {
				return fmt.Errorf("uhandled value type: %v", scalarValue)
			}
			fieldIndex += 1
		}
		// try to process section headers
		if err = parser.ProcessLine(sectionDataHeaders); err != nil {
			return fmt.Errorf("error processing section headers: %v", err)
		}
		// finally, section data.
		if err = parser.ProcessLine(sectionDataValues); err != nil {
			return fmt.Errorf("error processing section values: %v", err)
		}
	}

	return err
}
