package invoices

import (
	"ksef/cmd/ksef/commands/client"
	v2 "ksef/internal/client/v2"
	"ksef/internal/invoicesdb"
	invoicesdbconfig "ksef/internal/invoicesdb/config"
	inputprocessors "ksef/internal/sei/input_processors"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var importCommand = &cobra.Command{
	Use:   "import [input]",
	Short: "importuj faktury z pliku do bazy",
	Args:  cobra.ExactArgs(1),
	RunE:  importRun,
}

func init() {
	invoicesdbconfig.ImportFlags(importCommand.Flags())
	inputprocessors.GeneratorFlags(importCommand.Flags())
	importCommand.Flags().SortFlags = false
	InvoicesCommand.AddCommand(importCommand)
}

func importRun(cmd *cobra.Command, args []string) error {
	vip := viper.GetViper()
	ksefClient, err := client.InitClient(cmd, v2.WithoutTokenManager())
	if err != nil {
		return err
	}
	defer ksefClient.Close()
	// initialize the invoicesdb
	invoicesDB, err := invoicesdb.NewInvoicesDB(vip, invoicesdb.WithKSeFClient(ksefClient))
	if err != nil {
		return err
	}
	return invoicesDB.Import(
		cmd.Context(),
		vip,
		args[0],
		vip.GetBool(flagNameConfirm),
	)
}
