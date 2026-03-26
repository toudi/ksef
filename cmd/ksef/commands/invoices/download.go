package invoices

import (
	"ksef/cmd/ksef/commands/client"
	"ksef/cmd/ksef/flags"
	"ksef/internal/invoicesdb"
	downloaderconfig "ksef/internal/invoicesdb/downloader/config"
	kr "ksef/internal/keyring"
	"ksef/internal/logging"
	"ksef/internal/runtime"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	flagNameWorkersLong  = "workers"
	flagNameWorkersShort = "w"
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
	flagSet.IntP(flagNameWorkersLong, flagNameWorkersShort, 0, "Ilość workerów (domyślnie 0; wartość > 0 oznacza ilość współbieżnych wątków pobierających faktury dla wszystkich zarejestrowanych numerów NIP)")

	InvoicesCommand.AddCommand(downloadCommand)
}

func downloadRun(cmd *cobra.Command, _ []string) error {
	vip := viper.GetViper()
	workers := vip.GetInt(flagNameWorkersLong)
	if workers > 0 {
		return downloadRunParalell(cmd, vip, workers)
	}
	if err := runtime.CheckNIPIsSet(vip); err != nil {
		return err
	}

	nip, err := runtime.GetNIP(vip)
	if err != nil {
		return err
	}

	keyring, err := kr.NewKeyring(vip)
	if err != nil {
		logging.SeiLogger.Error("błąd inicjalizacji keyringu", "err", err)
		return err
	}
	defer keyring.Close()

	ksefClient, err := client.InitClient(cmd, vip, keyring)
	if err != nil {
		return err
	}
	defer ksefClient.Close()

	invoicesDB, err := invoicesdb.OpenForNIP(nip, vip, invoicesdb.WithKSeFClient(ksefClient))
	if err != nil {
		return err
	}

	downloaderConfig, err := downloaderconfig.GetDownloaderConfig(vip, "")
	if err != nil {
		return err
	}

	return invoicesDB.DownloadInvoices(
		cmd.Context(), vip, downloaderConfig,
		logging.DownloadLogger.With("nip", nip),
	)
}

func cloneViper(src *viper.Viper) *viper.Viper {
	newViper := viper.New()
	for _, key := range src.AllKeys() {
		newViper.Set(key, src.Get(key))
	}

	return newViper
}
