package jpk

import (
	"ksef/internal/invoicesdb/config"
	"ksef/internal/invoicesdb/jpk/generators"
	"ksef/internal/invoicesdb/jpk/generators/interfaces"
	"ksef/internal/invoicesdb/jpk/manager"
	monthlyregistry "ksef/internal/invoicesdb/monthly-registry"
	"ksef/internal/runtime"
	"os"
	"path/filepath"
	"slices"
	"time"

	"github.com/beevik/etree"
	"github.com/spf13/viper"
)

type JPK struct {
	manager   *manager.JPKManager
	registry  *monthlyregistry.Registry
	generator interfaces.JPKGenerator
	vip       *viper.Viper
	path      string
}

func NewJPK(month time.Time, vip *viper.Viper) (*JPK, error) {
	monthlyRegistry, err := monthlyregistry.OpenForMonth(vip, month)
	if err != nil {
		return nil, err
	}
	invoicesDBConfig := config.GetInvoicesDBConfig(vip)
	nip, _ := runtime.GetNIP(vip)

	path := filepath.Join(
		invoicesDBConfig.Root,
		runtime.GetEnvironmentId(vip),
		nip,
		month.Format("2006"),
		month.Format("01"),
	)

	manager, err := manager.Manager(vip, manager.WithMonthlyRegistry(monthlyRegistry))
	if err != nil {
		return nil, err
	}

	generator := generators.GetJPKGenerator(manager, month)

	return &JPK{
		registry:  monthlyRegistry,
		manager:   manager,
		vip:       vip,
		path:      path,
		generator: generator,
	}, nil
}

func (j *JPK) Generate() error {
	invoiceTypes := []monthlyregistry.InvoiceType{
		monthlyregistry.InvoiceTypeIssued,
		monthlyregistry.InvoiceTypeReceived,
	}

	for _, invoice := range j.registry.Invoices {
		if !slices.Contains(invoiceTypes, invoice.Type) {
			continue
		}

		fileName := j.registry.InvoiceFilename(invoice).XML

		if err := j.processInvoiceFile(fileName, invoice); err != nil {
			return err
		}
	}

	if document, err := j.generator.Document(); err != nil {
		return err
	} else {
		return j.writeToFile(
			document,
			filepath.Join(
				j.path,
				"jpk",
			),
		)
	}
}

func (j *JPK) processInvoiceFile(fileName string, invoice *monthlyregistry.Invoice) error {
	invoiceFile, err := os.Open(fileName)
	if err != nil {
		return err
	}

	defer invoiceFile.Close()

	xmlInvoice := etree.NewDocument()
	if _, err = xmlInvoice.ReadFrom(invoiceFile); err != nil {
		return err
	}

	return j.generator.ProcessInvoice(invoice, xmlInvoice)
}
