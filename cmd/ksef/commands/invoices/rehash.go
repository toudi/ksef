package invoices

import (
	"ksef/cmd/ksef/flags"
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

	// TODO: repopulate ksef ref no's

	return registry.Save()
}
