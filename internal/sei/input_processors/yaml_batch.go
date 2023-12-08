package inputprocessors

import (
	"fmt"
	"ksef/internal/sei/parser"
)

func (y *YAMLDecoder) processBatchSource(source map[string]interface{}, parser *parser.Parser) error {
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
