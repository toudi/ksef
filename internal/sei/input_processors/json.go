package inputprocessors

import (
	"encoding/json"
	"fmt"
	"ksef/internal/sei/parser"
	"os"
)

type JSONDecoder struct {
}

func JSONDecoder_Init() *JSONDecoder {
	return &JSONDecoder{}
}

func (j *JSONDecoder) Process(sourceFile string, parser *parser.Parser) error {
	var err error
	var serializedInvoice map[string]interface{}

	file, err := os.Open(sourceFile)
	if err != nil {
		return fmt.Errorf("unable to open source file: %v", err)
	}
	defer file.Close()

	if err = json.NewDecoder(file).Decode(&serializedInvoice); err != nil {
		return fmt.Errorf("error decoding JSON file. syntax error ?: %v", err)
	}

	// let's check if this is a bulk invoices array or a single one.
	if _, exists := serializedInvoice["common"]; exists {
		return j.processBatchSource(serializedInvoice, parser)
	} else {
		return j.processSingleInvoiceSource(serializedInvoice, parser)
	}
}

func (j *JSONDecoder) processSingleInvoiceSource(invoice map[string]interface{}, parser *parser.Parser) error {
	if err := processRecurse("", invoice, parser); err != nil {
		return fmt.Errorf("error running processRecurse: %v", err)
	}
	return parser.InvoiceReady()
}

func (j *JSONDecoder) processBatchSource(source map[string]interface{}, parser *parser.Parser) error {
	commonInvoiceData := source["common"]
	invoices := source["invoices"].([]interface{})
	var err error

	for _, invoice := range invoices {
		if err = processRecurse("", commonInvoiceData, parser); err != nil {
			return fmt.Errorf("unable to process common invoice data: %v", err)
		}
		if err = processRecurse("", invoice, parser); err != nil {
			return fmt.Errorf("unable to process invoice data: %v", err)
		}
		if err = parser.InvoiceReady(); err != nil {
			return fmt.Errorf("unable to call parser.InvoiceReady(): %v", err)
		}
	}

	return nil
}
