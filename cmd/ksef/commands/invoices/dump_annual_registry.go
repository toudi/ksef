package invoices

import (
	"ksef/cmd/ksef/flags"
	"ksef/core/invoicedb"
	"ksef/core/invoicedb/fs/metadata/registry"
	"ksef/core/invoicedb/fs/metadata/registry/monthly"
	"ksef/internal/resources"
	"ksef/internal/runtime"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	flagOutputFile        = "output"
	defaultOutputFileName = "registry-{{year}}.yaml"
)

var annualRegistryCommand = &cobra.Command{
	Use:   "annual-registry",
	Short: "Scala rejestry miesięczne z każdego istniejącego miesiąca bieżącego roku",
	RunE:  annualRegistryRun,
}

type annualInvoiceInfo struct {
	Issuer    string `yaml:"issuer,omitempty"`
	RefNo     string `yaml:"ref-no"`
	KSeFRefNo string `yaml:"ksef-ref-no"`
}

func init() {
	flagSet := annualRegistryCommand.Flags()
	flags.NIP(flagSet)
	flagSet.StringP(
		flagOutputFile,
		"o",
		defaultOutputFileName,
		"Szablon nazwy pliku wyjściowego. {{year}} Wstawia numer roku.",
	)
	InvoicesCommand.AddCommand(annualRegistryCommand)
}

func annualRegistryRun(cmd *cobra.Command, _ []string) error {
	vip := viper.GetViper()
	if err := runtime.CheckNIPIsSet(vip); err != nil {
		return err
	}

	var invoices []annualInvoiceInfo

	today := time.Now().Local()
	month := time.Date(today.Year(), 1, 1, 0, 0, 0, 0, time.Local)

	invoicesRoot, err := resources.InvoiceStorageRoot(vip)
	if err != nil {
		return err
	}

	// so basically all we have to do is to simply iterate over the months,
	// open the monthly registries and append the invoices which have been
	// synced with KSeF (not all of them will be) to the final array.
	for month.Before(today.AddDate(0, 1, 0)) {
		// there are portions of the code that need to be backend/implementation agnostic but this is not one of them.
		// this is a command that's used for managing the CLI (i.e. not a library) therefore it's safe to
		// forcefully initialize the filesystem-based API
		registry, err := registry.Open(invoicesRoot, invoicedb.YearMonth{Year: month.Year(), Month: int(month.Month())})
		// registry that does not exist is not a problem per se
		if err != nil && !os.IsNotExist(err) {
			return err
		}
		if err == nil && registry != nil {
			for _, invoice := range registry.Invoices() {
				if invoice.KSeFRefNo == "" {
					continue
				}
				issuerNIP := "" // for own invoices
				if invoice.Type > monthly.InvoiceTypeIssued {
					issuerNIP = invoice.Issuer.NIP
				}
				invoices = append(invoices, annualInvoiceInfo{
					RefNo:     invoice.RefNo,
					KSeFRefNo: invoice.KSeFRefNo,
					Issuer:    issuerNIP,
				})
			}
			if err = registry.Close(); err != nil {
				return err
			}
		}
		month = month.AddDate(0, 1, 0)
	}

	outputFileName := strings.NewReplacer(
		"{{year}}", strconv.Itoa(today.Year()),
	).Replace(vip.GetString(flagOutputFile))

	outputFile, err := os.Create(outputFileName)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	return yaml.NewEncoder(outputFile).Encode(invoices)
}
