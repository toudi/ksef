package inputprocessors

import (
	"fmt"
	"ksef/internal/sei/parser"
	"os"

	"github.com/goccy/go-yaml"
)

type YAMLDecoder struct{}

func YAMLDecoder_Init() *YAMLDecoder {
	return &YAMLDecoder{}
}

func (y *YAMLDecoder) Process(sourceFile string, parser *parser.Parser) error {
	var err error
	var serializedInvoice map[string]any

	file, err := os.Open(sourceFile)
	if err != nil {
		return fmt.Errorf("unable to open source file: %v", err)
	}
	defer file.Close()

	if err = yaml.NewDecoder(file).Decode(&serializedInvoice); err != nil {
		return fmt.Errorf("error decoding YAML file. syntax error ?: %v", err)
	}

	// let's check if this is a bulk invoices array or a single one.
	if _, exists := serializedInvoice["common"]; exists {
		return y.processBatchSource(serializedInvoice, parser)
	} else {
		return y.processSingleInvoiceSource(serializedInvoice, parser)
	}
}

func (y *YAMLDecoder) processSingleInvoiceSource(invoice map[string]any, parser *parser.Parser) error {
	if err := processRecurse("", invoice, parser); err != nil {
		return fmt.Errorf("error running processRecurse: %v", err)
	}
	return parser.InvoiceReady()
}

func (y *YAMLDecoder) processBatchSource(source map[string]any, parser *parser.Parser) error {
	commonInvoiceData := source["common"]
	invoices := source["invoices"].([]any)
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
