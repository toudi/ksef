package invoices

import (
	"ksef/cmd/ksef/commands/client"
	"ksef/cmd/ksef/flags"
	"ksef/internal/invoicesdb"
	downloaderconfig "ksef/internal/invoicesdb/downloader/config"
	"ksef/internal/runtime"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var downloadCommand = &cobra.Command{
	Use:   "download",
	Short: "pobierz nowe faktury z KSeF",
	RunE:  downloadRun,
}

func init() {
	flagSet := downloadCommand.Flags()
	flags.NIP(flagSet)
	downloaderconfig.DownloaderFlags(flagSet, "")
	runtime.CertProfileFlag(flagSet)

	InvoicesCommand.AddCommand(downloadCommand)
}

func downloadRun(cmd *cobra.Command, _ []string) error {
	vip := viper.GetViper()
	if err := runtime.CheckNIPIsSet(vip); err != nil {
		return err
	}

	nip, err := runtime.GetNIP(vip)
	if err != nil {
		return err
	}

	ksefClient, err := client.InitClient(cmd)
	if err != nil {
		return err
	}
	defer ksefClient.Close()

	invoicesDB, err := invoicesdb.OpenForNIP(
		nip, vip,
		invoicesdb.WithKSeFClient(ksefClient),
	)
	if err != nil {
		return err
	}

	downloaderConfig, err := downloaderconfig.GetDownloaderConfig(vip, "")
	if err != nil {
		return err
	}

	return invoicesDB.DownloadInvoices(cmd.Context(), vip, downloaderConfig)
}
