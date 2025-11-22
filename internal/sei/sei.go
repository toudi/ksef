package sei

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"ksef/internal/certsdb"
	"ksef/internal/interfaces"
	"ksef/internal/invoice"
	"ksef/internal/runtime"
	"ksef/internal/sei/generators"
	inputprocessors "ksef/internal/sei/input_processors"
	"ksef/internal/sei/parser"
	"ksef/internal/xml"
	"path/filepath"
	"strings"
)

var errUnknownSourceExtension = errors.New("unknown file extension")

// SEI is an government acronym for Structured Electronic Invoice
type SEI struct {
	env runtime.Gateway
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

func SEI_Init(env runtime.Gateway, conversionParams inputprocessors.InputProcessorConfig, initializers ...func(s *SEI)) (*SEI, error) {
	var err error

	var generator interfaces.Generator

	if generator, err = generators.Generator(conversionParams.Generator); err != nil {
		return nil, err
	}

	sei := &SEI{
		env:                  env,
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

// when the invoice is ready, we use the generator to convert it to
// xml.Node and persist to disk
// func (s *SEI) AddInvoice(invoice *invoice.Invoice) error {

// 	if s.registry.ContainsHash(invoiceHash) {
// 		logging.GenerateLogger.Info("faktura o obliczonej sumie kontrolnej już istnieje w zbiorze rejestru. no-op.")
// 		return nil
// 	}

// 	if s.IssuerTIN == "" {
// 		s.IssuerTIN = s.generator.IssuerTIN()
// 		s.registry.Issuer = s.IssuerTIN
// 	}

// 	if err := s.saveInvoice(); err != nil {
// 		return err
// 	}

// 	var invoiceMeta = invoices.InvoiceMetadata{
// 		Metadata:      invoice.Meta,
// 		InvoiceNumber: invoice.Number,
// 		IssueDate:     invoice.Issued.Format("2006-01-02"),
// 		Seller: invoices.InvoiceSubjectMetadata{
// 			NIP: s.IssuerTIN,
// 		},
// 		Offline: s.conversionParameters.OfflineMode,
// 	}

// 	certificate, err := s.getOfflineModeCertificate()
// 	if err != nil {
// 		return err
// 	}

// 	if err = s.registry.AddInvoice(invoiceMeta, invoiceHash, certificate); err != nil {
// 		return err
// 	}

// 	s.numInvoices += 1

// 	return err
// }

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
	s.parser.InvoiceReadyFunc = s.invoiceReady
	s.parser.LineHandler = s.generator.LineHandler

	processErr := processor.Process(sourceFile, s.parser)
	if processErr != nil {
		return processErr
	}
	return nil
	// return s.registry.Save("")
}

// func (s *SEI) saveInvoice() (err error) {
// 	fileName := filepath.Join(s.outputPath, fmt.Sprintf("invoice-%d.xml", s.numInvoices))
// 	destFile, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR, 0644)
// 	if err != nil {
// 		return fmt.Errorf("cannot create target file: %v", err)
// 	}
// 	if err = destFile.Truncate(0); err != nil {
// 		return err
// 	}
// 	defer destFile.Close()
// 	_, err = io.Copy(destFile, &s.invoiceContentBuffer)
// 	return err
// }

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

func (s *SEI) GetOfflineModeCertificate(i *invoice.Invoice) (*certsdb.Certificate, error) {
	if s.certificateForOfflineMode == nil {
		certsDB, err := certsdb.OpenOrCreate(s.env)
		if err != nil {
			return nil, err
		}
		offlineCert, err := certsDB.GetByUsage(
			certsdb.UsageOffline, i.IssuerNIP,
		)
		if err != nil {
			return nil, err
		}
		s.certificateForOfflineMode = &offlineCert
	}

	return s.certificateForOfflineMode, nil
}
