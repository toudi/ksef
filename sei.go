package ksef

import (
	"errors"
	"fmt"
	"ksef/common"
	inputprocessors "ksef/common/input_processors"
	"ksef/common/xml"
	"ksef/generators"
	"os"
	"path/filepath"
	"strings"
)

var errUnknownSourceExtension = errors.New("unknown file extension")

// SEI is an government acronym for Structured Electronic Invoice
type SEI struct {
	// conversion parameters for the input processor
	conversionParameters inputprocessors.InputProcessorConfig
	// xml files will be saved here
	outputPath string
	// how many invoices will be produced
	numInvoices int
	parser      *common.Parser
	IssuerTIN   string
	generator   common.Generator
}

func SEI_Init(outputPath string, conversionParams inputprocessors.InputProcessorConfig) (*SEI, error) {
	var err error

	if _, err = os.Stat(outputPath); os.IsNotExist(err) {
		// output path does not exist, let's try to create it
		if err = os.MkdirAll(outputPath, 0755); err != nil {
			return nil, fmt.Errorf("error creating dir: %v", err)
		}
	}
	var generator common.Generator

	if generator, err = generators.Generator(conversionParams.Generator); err != nil {
		return nil, err
	}

	return &SEI{
		outputPath:           outputPath,
		conversionParameters: conversionParams,
		parser:               &common.Parser{},
		generator:            generator,
	}, nil
}

// when the invoice is ready, we use the generator to convert it to
// xml.Node and persist to disk
func (s *SEI) AddInvoice(invoice *common.Invoice) error {
	invoiceXML, err := s.generator.InvoiceToXMLTree(invoice)
	if err != nil {
		return err
	}

	err = s.saveInvoice(invoiceXML)
	if err != nil {
		return err
	}

	s.numInvoices += 1

	if s.IssuerTIN == "" {
		s.IssuerTIN = s.generator.IssuerTIN()
	}

	return err
}

// this function loops over the entries within csv/xslx file and yields *invoice.Invoice objects
// these are then passed to AddInvoice function that saves them to disk.
func (s *SEI) ProcessSourceFile(sourceFile string) error {
	extension := strings.ToLower(filepath.Ext(sourceFile))
	var err error
	var processor inputprocessors.InputProcessor = nil

	if extension == ".csv" {
		processor, err = inputprocessors.CSVDecoder_Init(s.conversionParameters.CSV)
		if err != nil {
			return err
		}

	} else if extension == ".xlsx" {
		processor = inputprocessors.XLSXDecoder_Init(s.conversionParameters.XLSX)
	} else if extension == ".yaml" {
		processor = inputprocessors.YAMLDecoder_Init()
	}
	if processor == nil {
		return errUnknownSourceExtension
	}

	// set hook function to invoice ready signal
	s.parser.InvoiceReadyFunc = s.AddInvoice
	s.parser.LineHandler = s.generator.LineHandler

	return processor.Process(sourceFile, s.parser)
}

func (s *SEI) saveInvoice(root *xml.Node) error {
	destFileName := strings.Join([]string{s.outputPath, fmt.Sprintf("invoice-%d.xml", s.numInvoices)}, string(os.PathSeparator))
	destFile, err := os.OpenFile(destFileName, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("cannot create target file: %v", err)
	}
	if err = destFile.Truncate(0); err != nil {
		return err
	}
	defer destFile.Close()
	if _, err = destFile.WriteString("<?xml version=\"1.0\" encoding=\"utf-8\"?>\n"); err != nil {
		return err
	}
	return root.DumpToFile(destFile, 0)
}
