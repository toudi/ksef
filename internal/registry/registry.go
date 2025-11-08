package registry

import (
	"errors"
	"fmt"
	"ksef/internal/environment"
	"os"
	"path"
	"time"

	"gopkg.in/yaml.v3"
)

type DateType string

const (
	DateTypeIssue     DateType = "Issue"
	DateTypeInvoicing DateType = "Invoicing"
	DateTypeStorage   DateType = "PermanentStorage"
)

const registryName = "registry.yaml"

type InvoiceUploadResult struct {
	Filename string `yaml:"filename"`
	Checksum string `yaml:"checksum"`
	SeiRefNo string `yaml:"seiRefNo"`
	Failed   bool   `yaml:"failed,omitempty"` // whether the invoice was processed successfuly
}

type UploadSessionStatus struct {
	Processed bool                   `yaml:"processed"`
	Invoices  []*InvoiceUploadResult `yaml:"invoices"`
	UPO       []string               `yaml:"upo,omitempty"`
}

type QueryCriteria struct {
	DateType    DateType  `yaml:"invoicingDateType"`
	DateFrom    time.Time `yaml:"invoicingDateFrom"`
	DateTo      time.Time `yaml:"invoicingDateTo,omitempty"`
	SubjectType string    `yaml:"subjectType,omitempty"`
	Type        string    `yaml:"type,omitempty"`
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
	FullName                    string `yaml:"fullName,omitempty"`
}

type InvoiceQRCodes struct {
	Invoice     string `yaml:"invoice,omitempty"`
	Certificate string `yaml:"certificate,omitempty"`
}

type Invoice struct {
	ReferenceNumber     string         `yaml:"referenceNumber,omitempty"`
	KSeFReferenceNumber string         `yaml:"ksefReferenceNumber,omitempty"`
	QRCodes             InvoiceQRCodes `yaml:"qrcodes"`
	IssueDate           string         `yaml:"issueDate,omitempty"`
	SubjectFrom         InvoiceSubject `yaml:"subjectFrom,omitempty"`
	SubjectTo           InvoiceSubject `yaml:"subjectTo,omitempty"`
	InvoiceType         string         `yaml:"invoiceType,omitempty"`
	Net                 string         `yaml:"net,omitempty"`
	Vat                 string         `yaml:"vat,omitempty"`
	Gross               string         `yaml:"gross,omitempty"`
	Checksum            string         `yaml:"checksum,omitempty"`
	Offline             bool           `yaml:"offline,omitempty"`
}

type InvoiceRefId struct {
	ReferenceNumber     string `yaml:"invoiceRefNo"`
	KSeFReferenceNumber string `yaml:"ksefInvoiceRefNo"`
}

type PaymentId struct {
	SEIPaymentRefNo string   `yaml:"ksefPaymentRefNo"`
	InvoiceIDS      []string `yaml:"ksefInvoiceRefNumbers"`
}

type InvoiceRegistry struct {
	QueryCriteria  QueryCriteria                   `yaml:"queryCriteria,omitempty"`
	Environment    environment.Environment         `yaml:"environment,omitempty"`
	Invoices       []Invoice                       `yaml:"invoices,omitempty"`
	Issuer         string                          `yaml:"issuer,omitempty"`
	seiRefNoIndex  map[string]int                  `yaml:"-"`
	refNoIndex     map[string]int                  `yaml:"-"`
	checksumIndex  map[string]int                  `yaml:"-"`
	PaymentIds     []PaymentId                     `yaml:"payment-ids,omitempty"`
	UploadSessions map[string]*UploadSessionStatus `yaml:"upload-sessions,omitempty"` // map between upload session ID and list of seiRefNumbers
	sourcePath     string                          `yaml:"-"`
	Dir            string                          `yaml:"-"` // diretory for invoice registry. we cache it so that the registry knows where to save itself
	collection     *InvoiceCollection              `yaml:"-"` // cache invoice collection ptr
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

func LoadRegistry(dirName string) (*InvoiceRegistry, error) {
	var fileName = path.Join(dirName, registryName)
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
		registry.seiRefNoIndex[invoice.KSeFReferenceNumber] = index
		registry.refNoIndex[invoice.ReferenceNumber] = index
		if invoice.Checksum != "" {
			registry.checksumIndex[invoice.Checksum] = index
		}
	}
	registry.sourcePath = fileName
	registry.Dir = dirName
	return &registry, nil
}

func OpenOrCreate(dirName string) (*InvoiceRegistry, error) {
	registry, err := LoadRegistry(dirName)
	if registry == nil && err == ErrDoesNotExist {
		registry = NewRegistry()
		registry.sourcePath = path.Join(dirName, registryName)
		registry.Dir = dirName
		return registry, nil
	}

	if err != nil {
		return nil, fmt.Errorf("Unexpected error during opening registry file: %v", err)
	}

	return registry, nil
}
