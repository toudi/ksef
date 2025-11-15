package sei

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"ksef/internal/certsdb"
	"ksef/internal/client/v2/types/invoices"
	"ksef/internal/interfaces"
	"ksef/internal/invoice"
	"ksef/internal/logging"
	"ksef/internal/registry"
	"ksef/internal/runtime"
	"ksef/internal/sei/generators"
	inputprocessors "ksef/internal/sei/input_processors"
	"ksef/internal/sei/parser"
	"ksef/internal/xml"
	"os"
	"path/filepath"
	"strings"
)

var errUnknownSourceExtension = errors.New("unknown file extension")

// SEI is an government acronym for Structured Electronic Invoice
type SEI struct {
	env runtime.Gateway
	// conversion parameters for the input processor
	conversionParameters inputprocessors.InputProcessorConfig
	// xml files will be saved here
	outputPath string
	// how many invoices will be produced
	numInvoices               int
	parser                    *parser.Parser
	IssuerTIN                 string
	generator                 interfaces.Generator
	registry                  *registry.InvoiceRegistry
	invoiceContentBuffer      bytes.Buffer
	certificateForOfflineMode *certsdb.Certificate
}

func SEI_Init(env runtime.Gateway, outputPath string, conversionParams inputprocessors.InputProcessorConfig) (*SEI, error) {
	var err error

	var r *registry.InvoiceRegistry
	if r, err = registry.OpenOrCreate(outputPath); err != nil {
		return nil, err
	}
	if r.Environment == "" {
		r.Environment = env
	}

	var generator interfaces.Generator

	if generator, err = generators.Generator(conversionParams.Generator); err != nil {
		return nil, err
	}

	return &SEI{
		env:                  env,
		outputPath:           outputPath,
		conversionParameters: conversionParams,
		parser:               &parser.Parser{},
		generator:            generator,
		registry:             r,
		numInvoices:          len(r.Invoices),
	}, nil
}

// when the invoice is ready, we use the generator to convert it to
// xml.Node and persist to disk
func (s *SEI) AddInvoice(invoice *invoice.Invoice) error {
	s.invoiceContentBuffer.Reset()

	invoiceXML, err := s.generator.InvoiceToXMLTree(invoice)
	if err != nil {
		return err
	}

	invoiceHash, err := s.calculateInvoiceHash(invoiceXML)
	if err != nil {
		return err
	}

	if s.registry.ContainsHash(invoiceHash) {
		logging.GenerateLogger.Info("faktura o obliczonej sumie kontrolnej już istnieje w zbiorze rejestru. no-op.")
		return nil
	}

	if s.IssuerTIN == "" {
		s.IssuerTIN = s.generator.IssuerTIN()
		s.registry.Issuer = s.IssuerTIN
	}

	if err := s.saveInvoice(); err != nil {
		return err
	}

	var invoiceMeta = invoices.InvoiceMetadata{
		InvoiceNumber: invoice.Number,
		IssueDate:     invoice.Issued.Format("2006-01-02"),
		Seller: invoices.InvoiceSubjectMetadata{
			NIP: s.IssuerTIN,
		},
		Offline: s.conversionParameters.OfflineMode,
	}

	certificate, err := s.getOfflineModeCertificate()
	if err != nil {
		return err
	}

	if err = s.registry.AddInvoice(invoiceMeta, invoiceHash, certificate); err != nil {
		return err
	}

	s.numInvoices += 1

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
	} else if extension == ".json" {
		processor = inputprocessors.JSONDecoder_Init()
	}
	if processor == nil {
		return errUnknownSourceExtension
	}

	// set hook function to invoice ready signal
	s.parser.InvoiceReadyFunc = s.AddInvoice
	s.parser.LineHandler = s.generator.LineHandler

	processErr := processor.Process(sourceFile, s.parser)
	if processErr != nil {
		return processErr
	}
	return s.registry.Save("")
}

func (s *SEI) saveInvoice() (err error) {
	fileName := filepath.Join(s.outputPath, fmt.Sprintf("invoice-%d.xml", s.numInvoices))
	destFile, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("cannot create target file: %v", err)
	}
	if err = destFile.Truncate(0); err != nil {
		return err
	}
	defer destFile.Close()
	_, err = io.Copy(destFile, &s.invoiceContentBuffer)
	return err
}

func (s *SEI) calculateInvoiceHash(root *xml.Node) (checksum string, err error) {
	if _, err = s.invoiceContentBuffer.WriteString("<?xml version=\"1.0\" encoding=\"utf-8\"?>\n"); err != nil {
		return checksum, err
	}
	if err = root.DumpToWriter(&s.invoiceContentBuffer, 0); err != nil {
		return checksum, err
	}
	var hashBytes = sha256.Sum256(s.invoiceContentBuffer.Bytes())
	checksum = hex.EncodeToString(hashBytes[:])
	return checksum, nil
}

func (s *SEI) getOfflineModeCertificate() (*certsdb.Certificate, error) {
	if s.certificateForOfflineMode == nil {
		certsDB, err := certsdb.OpenOrCreate(s.env)
		if err != nil {
			return nil, err
		}
		offlineCert, err := certsDB.GetByUsage(
			certsdb.UsageOffline, s.IssuerTIN,
		)
		if err != nil {
			return nil, err
		}
		s.certificateForOfflineMode = &offlineCert
	}

	return s.certificateForOfflineMode, nil
}
