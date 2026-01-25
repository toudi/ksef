package jpk

import (
	"ksef/internal/invoicesdb/config"
	"ksef/internal/invoicesdb/jpk/types"
	monthlyregistry "ksef/internal/invoicesdb/monthly-registry"
	subjectsettings "ksef/internal/invoicesdb/subject-settings"
	"ksef/internal/money"
	"ksef/internal/runtime"
	"os"
	"path/filepath"
	"slices"
	"time"

	"github.com/beevik/etree"
	"github.com/spf13/viper"
)

type ControlRow struct {
	VAT money.MonetaryValue
}

type JPK struct {
	manager      *JPKManager
	registry     *monthlyregistry.Registry
	sjs          *subjectsettings.JPKSettings
	vip          *viper.Viper
	path         string
	Income       []*types.Invoice
	IncomeCtrl   *ControlRow
	Purchase     []*types.Invoice
	PurchaseCtrl *ControlRow
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

	ss, err := subjectsettings.OpenOrCreate(
		filepath.Join(
			invoicesDBConfig.Root,
			runtime.GetEnvironmentId(vip),
			nip,
		),
	)
	if err != nil {
		return nil, err
	}

	manager, err := Manager(vip, WithMonthlyRegistry(monthlyRegistry))
	if err != nil {
		return nil, err
	}

	return &JPK{
		registry:     monthlyRegistry,
		manager:      manager,
		vip:          vip,
		path:         path,
		sjs:          ss.JPK,
		IncomeCtrl:   &ControlRow{},
		PurchaseCtrl: &ControlRow{},
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

	return j.Save(
		filepath.Join(
			j.path,
			"jpk",
		),
	)
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

	if invoice.Type == monthlyregistry.InvoiceTypeIssued {
		if err = j.AddIncome(xmlInvoice, invoice); err != nil {
			return err
		}
	} else {
		if err = j.AddReceived(xmlInvoice, invoice); err != nil {
			return err
		}
	}

	return nil
}
