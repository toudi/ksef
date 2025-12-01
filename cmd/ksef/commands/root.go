package commands

import (
	"fmt"
	"ksef/cmd/ksef/commands/authorization"
	"ksef/cmd/ksef/commands/certificates"
	"ksef/cmd/ksef/commands/download"
	"ksef/internal/config"
	"ksef/internal/logging"
	"ksef/internal/runtime"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	flagTestGateway = "test-gateway"
	flagDemoGateway = "demo-gateway"
)

var RootCommand = &cobra.Command{
	Use:               "ksef",
	Long:              "KSeF - aplikacja kliencka",
	PersistentPreRunE: setContext,
}

var (
	configFile string
	logOutput  string
)

func init() {
	runtime.SetGateway(viper.GetViper(), runtime.ProdGateway)
	RootCommand.PersistentFlags().StringVarP(&logOutput, "log", "l", "-", "wyjście logowania (wartość - oznacza wyjście standardowe)")
	RootCommand.PersistentFlags().StringVarP(&configFile, "config", "c", "config.yaml", "lokalizacja pliku konfiguracyjnego")
	RootCommand.PersistentFlags().BoolFuncP("verbose", "v", "tryb verbose", func(s string) error {
		logging.Verbose = true
		return nil
	})
	RootCommand.PersistentFlags().BoolFuncP(flagTestGateway, "t", "Użyj bramki testowej", func(s string) error {
		runtime.SetGateway(viper.GetViper(), runtime.TestGateway)
		return nil
	})
	RootCommand.PersistentFlags().BoolFuncP(flagDemoGateway, "", "Użyj bramki przedprodukcyjnej (demo)", func(s string) error {
		runtime.SetGateway(viper.GetViper(), runtime.DemoGateway)
		return nil
	})
	config.DataDirFlag(RootCommand)
	RootCommand.PersistentFlags().SortFlags = false

	RootCommand.AddCommand(authorization.AuthCommand)
	RootCommand.AddCommand(certificates.CertificatesCommand)
	// RootCommand.AddCommand(syncInvoicesCommand)
	RootCommand.AddCommand(download.DownloadCommand)
	RootCommand.AddCommand(uploadCommand)
	RootCommand.AddCommand(statusCommand)
	RootCommand.AddCommand(renderPDFCommand)
	RootCommand.CompletionOptions = cobra.CompletionOptions{DisableDefaultCmd: true}
}

func setContext(cmd *cobra.Command, _ []string) error {
	var err error

	viper.SetConfigType("yaml")
	viper.SetEnvPrefix("ksef")
	viper.AutomaticEnv()

	if err = viper.BindPFlags(cmd.Flags()); err != nil {
		return err
	}

	if configFile != "" {
		viper.SetConfigFile(configFile)
		if err = viper.ReadInConfig(); err != nil && !os.IsNotExist(err) {
			return err
		}
	}

	if err = logging.InitLogging(logOutput); err != nil {
		fmt.Printf("[ ERROR ] Błąd inicjalizacji logowania: %v", err)
		return err
	}

	logging.SeiLogger.Info("start programu")
	logging.SeiLogger.Info("wybrane środowisko", "env", runtime.GetGateway(viper.GetViper()))

	return nil
}
