package jpk

import (
	"errors"
	"ksef/internal/invoicesdb/annotations"
	"ksef/internal/invoicesdb/config"
	"ksef/internal/invoicesdb/jpk/constants"
	"ksef/internal/invoicesdb/jpk/generators"
	"ksef/internal/invoicesdb/jpk/generators/interfaces"
	monthlyregistry "ksef/internal/invoicesdb/monthly-registry"
	subjectsettings "ksef/internal/invoicesdb/subject-settings"
	"ksef/internal/runtime"
	"os"
	"path/filepath"
	"time"

	"github.com/beevik/etree"
	"github.com/spf13/viper"
)

type JPK struct {
	manager   *annotations.Annotations
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

	manager, err := annotations.Manager(vip, annotations.WithMonthlyRegistry(monthlyRegistry))
	if err != nil {
		return nil, err
	}

	// Load subject settings directly - the path is derived from the registry
	settingsPath := filepath.Join(
		monthlyRegistry.Dir(),
		"..",
		"..",
	)
	subjectSettings, err := subjectsettings.OpenOrCreate(settingsPath)
	if err != nil {
		// subject settings are optional, so we just continue without them
		subjectSettings = &subjectsettings.SubjectSettings{
			JPK: &subjectsettings.JPKSettings{
				Surplus: subjectsettings.SurplusAction{
					CarryOver: true,
				},
			},
		}
	}

	if subjectSettings.JPK == nil {
		return nil, errors.New("brak ustawień JPK")
	}

	refundMode := vip.GetString(constants.FlagNameRefundMode)
	if refundMode != "" {
		subjectSettings.JPK.Surplus.CarryOver = false
		subjectSettings.JPK.Surplus.Refund = refundMode
	}

	offsetTax := vip.GetString(constants.FlagNameOffsetTaxCode)
	if offsetTax != "" && offsetTax != subjectSettings.JPK.Surplus.OffsetTax {
		subjectSettings.JPK.Surplus.CarryOver = false
		subjectSettings.JPK.Surplus.OffsetTax = offsetTax
	}

	if err := subjectSettings.JPK.Validate(); err != nil {
		return nil, err
	}

	generator := generators.GetJPKGenerator(manager, subjectSettings, month)

	return &JPK{
		registry:  monthlyRegistry,
		manager:   manager,
		vip:       vip,
		path:      path,
		generator: generator,
	}, nil
}

func (j *JPK) Generate() error {
	for _, invoice := range j.registry.JPKEligibleInvoices() {
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
