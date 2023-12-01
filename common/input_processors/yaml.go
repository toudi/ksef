package inputprocessors

import (
	"fmt"
	"ksef/common"
	"os"

	"gopkg.in/yaml.v3"
)

type YAMLDecoder struct {
}

func YAMLDecoder_Init() *YAMLDecoder {
	return &YAMLDecoder{}
}

func (y *YAMLDecoder) Process(sourceFile string, parser *common.Parser) error {
	var err error
	var serializedInvoice map[string]interface{}

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

func (y *YAMLDecoder) FeedLine() ([]string, error) {
	return []string{}, nil
}
