package registry

import (
	"errors"
	"fmt"
	"ksef/internal/config"
	"os"
	"path"
	"time"

	"gopkg.in/yaml.v3"
)

const registryName = "registry.yaml"

type InvoiceUploadResult struct {
	Filename string `yaml:"filename"`
	Checksum string `yaml:"checksum"`
	SeiRefNo string `yaml:"seiRefNo"`
	Failed   bool   `yaml:"failed,omitzero"` // whether the invoice was processed successfuly
}

type UploadSessionStatus struct {
	Processed bool                   `yaml:"processed"`
	Invoices  []*InvoiceUploadResult `yaml:"invoices"`
}

type QueryCriteria struct {
	DateFrom    time.Time `json:"invoicingDateFrom" yaml:"invoicingDateFrom"`
	DateTo      time.Time `json:"invoicingDateTo"   yaml:"invoicingDateTo"`
	SubjectType string    `json:"subjectType"       yaml:"subjectType"`
	Type        string    `json:"type"              yaml:"type"`
}

type invoiceSubjectQueryResponse struct {
	IssuedBy struct {
		TIN string `json:"identifier"`
	} `json:"issuedByIdentifier"`
	Issuer struct {
		FullName string `json:"fullName"`
	} `json:"issuedByName"`
}

type InvoiceSubject struct {
	invoiceSubjectQueryResponse `       yaml:"-"`
	TIN                         string `yaml:"NIP"`
	FullName                    string `yaml:"fullName"`
}

type Invoice struct {
	ReferenceNumber    string         `json:"invoiceReferenceNumber" yaml:"referenceNumber,omitempty"`
	SEIReferenceNumber string         `json:"ksefReferenceNumber"    yaml:"ksefReferenceNumber,omitempty"`
	SEIQRCode          string         `json:"-"                      yaml:"qrcode-url,omitempty"`
	InvoicingDate      string         `json:"invoicingDate"          yaml:"invoicingDate,omitempty"`
	SubjectFrom        InvoiceSubject `json:"subjectBy,omitempty"    yaml:"subjectFrom,omitempty"`
	SubjectTo          InvoiceSubject `json:"subjectTo,omitempty"    yaml:"subjectTo,omitempty"`
	InvoiceType        string         `json:"invoiceType"            yaml:"invoiceType,omitempty"`
	Net                string         `json:"net"                    yaml:"net,omitempty"`
	Vat                string         `json:"vat"                    yaml:"vat,omitempty"`
	Gross              string         `json:"gross"                  yaml:"gross,omitempty"`
	Checksum           string         `                              yaml:"checksum,omitempty"`
}

type InvoiceRefId struct {
	ReferenceNumber    string `json:"invoiceRefNo"     yaml:"invoiceRefNo"`
	SEIReferenceNumber string `json:"ksefInvoiceRefNo" yaml:"ksefInvoiceRefNo"`
}

type PaymentId struct {
	SEIPaymentRefNo string   `yaml:"ksefPaymentRefNo"`
	InvoiceIDS      []string `yaml:"ksefInvoiceRefNumbers"`
}

type InvoiceRegistry struct {
	QueryCriteria  QueryCriteria                   `json:"queryCriteria" yaml:"queryCriteria,omitempty"`
	Environment    config.APIEnvironment           `                     yaml:"environment"`
	Invoices       []Invoice                       `                     yaml:"invoices,omitempty"`
	Issuer         string                          `                     yaml:"issuer,omitempty"`
	seiRefNoIndex  map[string]int                  `                     yaml:"-"`
	refNoIndex     map[string]int                  `                     yaml:"-"`
	checksumIndex  map[string]int                  `                     yaml:"-"`
	PaymentIds     []PaymentId                     `                     yaml:"payment-ids,omitempty"`
	UploadSessions map[string]*UploadSessionStatus `yaml:"upload-sessions"` // map between upload session ID and list of seiRefNumbers
	sourcePath     string
	collection     *InvoiceCollection `yaml:"-"` // cache invoice collection ptr
}

var ErrDoesNotExist = errors.New("registry file does not exist")

func NewRegistry() *InvoiceRegistry {
	_registry := &InvoiceRegistry{
		seiRefNoIndex: make(map[string]int),
		refNoIndex:    make(map[string]int),
		checksumIndex: make(map[string]int),
	}

	return _registry
}

func (r *InvoiceRegistry) Save(fileName string) error {
	if fileName == "" {
		fileName = r.sourcePath
	}
	if fileName == "" {
		return fmt.Errorf("fileName not specified")
	}
	destPath := path.Dir(fileName)
	if err := os.MkdirAll(destPath, 0755); err != nil {
		return fmt.Errorf("unable to create registry dir: %v", err)
	}
	destFile, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("unable to create registry file: %v", err)
	}
	return yaml.NewEncoder(destFile).Encode(r)
}

func LoadRegistry(fileName string) (*InvoiceRegistry, error) {
	var registry InvoiceRegistry
	registry.seiRefNoIndex = make(map[string]int)
	registry.refNoIndex = make(map[string]int)
	registry.checksumIndex = make(map[string]int)
	reader, err := os.Open(fileName)
	if err != nil {
		return nil, ErrDoesNotExist
	}
	if err = yaml.NewDecoder(reader).Decode(&registry); err != nil {
		return nil, fmt.Errorf("unable to decode invoice registry: %v", err)
	}
	for index, invoice := range registry.Invoices {
		registry.seiRefNoIndex[invoice.SEIReferenceNumber] = index
		registry.refNoIndex[invoice.ReferenceNumber] = index
		if invoice.Checksum != "" {
			registry.checksumIndex[invoice.Checksum] = index
		}
	}
	registry.sourcePath = fileName
	return &registry, nil
}

func OpenOrCreate(dirName string) (*InvoiceRegistry, error) {
	registry, err := LoadRegistry(path.Join(dirName, registryName))
	if registry == nil && err == ErrDoesNotExist {
		registry = NewRegistry()
		return registry, nil
	}

	if err != nil {
		return nil, fmt.Errorf("Unexpected error during opening registry file: %v", err)
	}

	return registry, nil
}
