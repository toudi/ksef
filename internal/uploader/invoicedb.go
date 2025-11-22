package uploader

import (
	"io"
	"os"
	"path"
	"time"

	"go.yaml.in/yaml/v3"
)

const (
	dbname = "invoices.yaml"
)

type Invoice struct {
	RefNo          string    `yaml:"ref-no"`
	KSeFRefNo      string    `yaml:"ksef-ref-no"`
	Contents       string    `yaml:"contents"`
	Checksum       string    `yaml:"checksum"`
	GenerationTime time.Time `yaml:"generation-time"`
	Corrections    []Invoice `yaml:"corrections,omitempty"`
}

type InvoiceDB struct {
	Invoices            []Invoice
	filename            string
	invoiceByRefNoIndex map[string]int
}

func InvoiceDB_OpenOrCreate(dir string) (*InvoiceDB, error) {
	var db = &InvoiceDB{
		invoiceByRefNoIndex: make(map[string]int),
	}
	var filename = path.Join(dir, dbname)
	dbFile, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) || err == io.EOF {
			db.filename = filename
			return db, nil
		}
		return nil, err
	}
	defer dbFile.Close()
	var decoder = yaml.NewDecoder(dbFile)
	err = decoder.Decode(&db.Invoices)
	if err == io.EOF {
		err = nil
	}
	db.filename = filename
	if err == nil {
		db.loadInvoiceByRefNoIndex()
	}
	return db, err
}

func (idb *InvoiceDB) Save() error {
	dbFile, err := os.Create(idb.filename)
	if err != nil {
		return err
	}
	var encoder = yaml.NewEncoder(dbFile)
	return encoder.Encode(idb.Invoices)
}

func (idb *InvoiceDB) loadInvoiceByRefNoIndex() {
	for index, invoice := range idb.Invoices {
		idb.invoiceByRefNoIndex[invoice.RefNo] = index
	}
}
