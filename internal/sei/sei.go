package sei

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"ksef/internal/certsdb"
	"ksef/internal/interfaces"
	"ksef/internal/sei/generators"
	inputprocessors "ksef/internal/sei/input_processors"
	"ksef/internal/sei/parser"
	"ksef/internal/xml"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

var errUnknownSourceExtension = errors.New("unknown file extension")

// SEI is an government acronym for Structured Electronic Invoice
type SEI struct {
	// conversion parameters for the input processor
	conversionParameters inputprocessors.InputProcessorConfig
	// this function will be called each time the invoice will be ready
	invoiceReadyFunc          func(i *ParsedInvoice) error
	parser                    *parser.Parser
	IssuerTIN                 string
	generator                 interfaces.Generator
	invoiceContentBuffer      bytes.Buffer
	certificateForOfflineMode *certsdb.Certificate
	invoiceExistsFunc         func(refNo string) bool
}

func SEI_Init(vip *viper.Viper, initializers ...func(s *SEI)) (*SEI, error) {
	var err error
	conversionParams := inputprocessors.GetInputProcessorConfig(vip)

	var generator interfaces.Generator

	if generator, err = generators.Generator(conversionParams.Generator); err != nil {
		return nil, err
	}

	sei := &SEI{
		conversionParameters: conversionParams,
		parser:               &parser.Parser{},
		generator:            generator,
		invoiceExistsFunc:    func(refNo string) bool { return false },
	}

	for _, initializer := range initializers {
		initializer(sei)
	}

	return sei, nil
}

// this function loops over the entries within csv/xslx file and yields *invoice.Invoice objects
// these are then passed to AddInvoice function that saves them to disk.
func (s *SEI) ProcessSourceFile(sourceFile string) error {
	processor, err := s.getProcessorFromFilename(sourceFile)
	if err != nil {
		return err
	}

	// set hook function to invoice ready signal
	s.parser.InvoiceReadyFunc = s.invoiceReady
	s.parser.LineHandler = s.generator.LineHandler

	processErr := processor.Process(sourceFile, s.parser)
	if processErr != nil {
		return processErr
	}
	return nil
}

func (s *SEI) ProcessReader(src io.Reader, filename string) error {
	processor, err := s.getProcessorFromFilename(filename)
	if err != nil {
		return err
	}

	// set hook function to invoice ready signal
	s.parser.InvoiceReadyFunc = s.invoiceReady
	s.parser.LineHandler = s.generator.LineHandler

	processErr := processor.ProcessReader(src, s.parser)
	if processErr != nil {
		return processErr
	}
	return nil
}

func (s *SEI) getProcessorFromFilename(filename string) (processor inputprocessors.InputProcessor, err error) {
	extension := strings.ToLower(filepath.Ext(filename))

	if extension == ".csv" {
		processor, err = inputprocessors.CSVDecoder_Init(s.conversionParameters.CSV)
		if err != nil {
			return nil, err
		}

	} else if extension == ".xlsx" {
		processor = inputprocessors.XLSXDecoder_Init(s.conversionParameters.XLSX)
	} else if extension == ".yaml" {
		processor = inputprocessors.YAMLDecoder_Init()
	} else if extension == ".json" {
		processor = inputprocessors.JSONDecoder_Init()
	}
	if processor == nil {
		return nil, errUnknownSourceExtension
	}

	return processor, nil
}

func (s *SEI) calculateInvoiceHash(root *xml.Node) (checksum string, err error) {
	if _, err = s.invoiceContentBuffer.WriteString("<?xml version=\"1.0\" encoding=\"utf-8\"?>\n"); err != nil {
		return checksum, err
	}
	if err = root.DumpToWriter(&s.invoiceContentBuffer, 0); err != nil {
		return checksum, err
	}
	hashBytes := sha256.Sum256(s.invoiceContentBuffer.Bytes())
	checksum = hex.EncodeToString(hashBytes[:])
	return checksum, nil
}
