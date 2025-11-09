package registry

import (
	"errors"
	"fmt"
	"ksef/internal/config"
	"ksef/internal/registry/types"
	"os"
	"path"

	"gopkg.in/yaml.v3"
)

const registryName = "registry.yaml"

var ErrDoesNotExist = errors.New("registry file does not exist")

type InvoiceRefId struct {
	ReferenceNumber     string `yaml:"invoiceRefNo"`
	KSeFReferenceNumber string `yaml:"ksefInvoiceRefNo"`
}

type PaymentId struct {
	SEIPaymentRefNo string   `yaml:"ksefPaymentRefNo"`
	InvoiceIDS      []string `yaml:"ksefInvoiceRefNumbers"`
}

type InvoiceRegistry struct {
	Sync           types.SyncConfig                      `yaml:"sync,omitempty"`
	Environment    config.Gateway                        `yaml:"environment,omitempty"`
	Invoices       []types.Invoice                       `yaml:"invoices,omitempty"`
	Issuer         string                                `yaml:"issuer,omitempty"`
	seiRefNoIndex  map[string]int                        `yaml:"-"`
	refNoIndex     map[string]int                        `yaml:"-"`
	checksumIndex  map[string]int                        `yaml:"-"`
	PaymentIds     []PaymentId                           `yaml:"payment-ids,omitempty"`
	UploadSessions map[string]*types.UploadSessionStatus `yaml:"upload-sessions,omitempty"` // map between upload session ID and list of seiRefNumbers
	sourcePath     string                                `yaml:"-"`
	Dir            string                                `yaml:"-"` // diretory for invoice registry. we cache it so that the registry knows where to save itself
	collection     *InvoiceCollection                    `yaml:"-"` // cache invoice collection ptr
}

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
