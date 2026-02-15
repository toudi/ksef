package invoices

import (
	"ksef/cmd/ksef/flags"
	annualregistry "ksef/internal/invoicesdb/annual-registry"
	monthlyregistry "ksef/internal/invoicesdb/monthly-registry"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rehashCommand = &cobra.Command{
	Use:   "rehash",
	Short: "odczytaj i zapisz rejestr, uzupełniając brakujące pola",
	RunE:  rehashRegistryRun,
	Args:  cobra.ExactArgs(2),
}

func init() {
	flagSet := rehashCommand.Flags()
	flags.NIP(flagSet)
	rehashCommand.MarkFlagRequired(flags.FlagNameNIP)
	InvoicesCommand.AddCommand(rehashCommand)
}

func rehashRegistryRun(cmd *cobra.Command, args []string) error {
	vip := viper.GetViper()

	yearInt, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}

	monthInt, err := strconv.Atoi(args[1])
	if err != nil {
		return err
	}

	month := time.Date(
		yearInt,
		time.Month(monthInt),
		1,
		0, 0, 0, 0, time.Local,
	)

	registry, err := monthlyregistry.OpenForMonth(
		vip, month,
	)
	if err != nil {
		return err
	}

	annualRegistry, err := annualregistry.OpenForMonth(
		vip, month,
	)
	if err != nil {
		return err
	}

	for _, invoice := range registry.Invoices {
		if invoice.KSeFRefNo == "" {
			continue
		}

		if err = annualRegistry.UpdateInvoiceByChecksum(
			invoice.Checksum,
			func(_invoice *annualregistry.Invoice) error {
				_invoice.KSeFRefNo = invoice.KSeFRefNo
				return nil
			},
		); err != nil {
			return err
		}
	}

	if err = registry.Save(); err != nil {
		return err
	}

	return annualRegistry.Save()
}
