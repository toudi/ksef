package inputprocessors

import (
	"fmt"
	"ksef/common"
)

func (y *YAMLDecoder) processBatchSource(source map[string]interface{}, parser *common.Parser) error {
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
