package invoices

import (
	"errors"
	"ksef/cmd/ksef/commands/client"
	v2 "ksef/internal/client/v2"
	"ksef/internal/invoicesdb"
	invoicesdbconfig "ksef/internal/invoicesdb/config"
	statuscheckerconfig "ksef/internal/invoicesdb/status-checker/config"
	uploaderconfig "ksef/internal/invoicesdb/uploader/config"
	kr "ksef/internal/keyring"
	"ksef/internal/logging"
	"ksef/internal/runtime"
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

var (
	errClientInit     = errors.New("error during client initialization")
	errInvoicesDBInit = errors.New("error during invoicesDB init")
)

func init() {
	invoicesdbconfig.ImportFlags(importCommand.Flags())
	inputprocessors.GeneratorFlags(importCommand.Flags())
	runtime.CertProfileFlag(importCommand.Flags())
	uploaderconfig.UploaderFlags(importCommand.Flags())
	statuscheckerconfig.StatusCheckerFlags(importCommand.Flags())

	importCommand.Flags().SortFlags = false
	InvoicesCommand.AddCommand(importCommand)
}

func importRun(cmd *cobra.Command, args []string) error {
	vip := viper.GetViper()

	keyring, err := kr.NewKeyring(vip)
	if err != nil {
		logging.SeiLogger.Error("błąd inicjalizacji keyringu", "err", err)
		return err
	}
	defer keyring.Close()

	ksefClient, err := client.InitClient(cmd, vip, keyring, v2.WithoutTokenManager())
	if err != nil {
		return errors.Join(errClientInit, err)
	}
	defer ksefClient.Close()
	// initialize the invoicesdb
	invoicesDB, err := invoicesdb.NewInvoicesDB(
		vip,
		invoicesdb.WithKSeFClient(ksefClient),
		invoicesdb.WithoutInitializingPrefix(),
	)
	if err != nil {
		return errors.Join(errInvoicesDBInit, err)
	}
	return invoicesDB.Import(
		cmd.Context(),
		vip,
		args[0],
		vip.GetBool(flagNameConfirm),
	)
}
