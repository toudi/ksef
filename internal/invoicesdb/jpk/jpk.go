package jpk

import (
	"ksef/internal/invoicesdb/config"
	monthlyregistry "ksef/internal/invoicesdb/monthly-registry"
	"ksef/internal/money"
	"ksef/internal/runtime"
	"path/filepath"
	"slices"
	"time"

	"github.com/spf13/viper"
)

type Amounts struct {
	base map[Extractor]money.MonetaryValue
	vat  map[Extractor]money.MonetaryValue
}

type JPK struct {
	registry *monthlyregistry.Registry
	vip      *viper.Viper
	path     string
	amounts  *Amounts
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

	return &JPK{
		registry: monthlyRegistry,
		vip:      vip,
		path:     path,
		amounts: &Amounts{
			base: make(map[Extractor]money.MonetaryValue),
			vat:  make(map[Extractor]money.MonetaryValue),
		},
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

		fileName := j.registry.InvoiceFilename(invoice)

		xmlInvoice, _, err := monthlyregistry.ParseInvoice(fileName.XML)
		if err != nil {
			return err
		}

		if invoice.Type == monthlyregistry.InvoiceTypeIssued {
			if err = j.AddIncome(xmlInvoice); err != nil {
				return err
			}
		} else {
			if err = j.AddReceived(xmlInvoice); err != nil {
				return err
			}
		}
	}

	return j.Save(
		filepath.Join(
			j.path,
			"jpk",
			"jpk-v7m.xml",
		),
	)
}
