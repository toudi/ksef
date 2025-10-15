package commands

import (
	"context"
	"fmt"
	"ksef/cmd/ksef/commands/authorization"
	"ksef/cmd/ksef/commands/certificates"
	appCtx "ksef/cmd/ksef/context"
	"ksef/internal/config"
	environmentPkg "ksef/internal/environment"
	"ksef/internal/logging"

	"github.com/spf13/cobra"
)

var RootCommand = &cobra.Command{
	Use:               "ksef",
	Long:              "KSeF - aplikacja kliencka",
	PersistentPreRunE: setContext,
}

var (
	env        environmentPkg.Environment = environmentPkg.Production
	configFile string
	logOutput  string
)

func init() {
	RootCommand.PersistentFlags().StringVarP(&logOutput, "log", "l", "-", "wyjście logowania (wartość - oznacza wyjście standardowe)")
	RootCommand.PersistentFlags().StringVarP(&configFile, "config", "c", "", "lokalizacja pliku konfiguracyjnego")
	RootCommand.PersistentFlags().BoolFuncP("verbose", "v", "tryb verbose", func(s string) error {
		logging.Verbose = true
		return nil
	})
	RootCommand.PersistentFlags().BoolFuncP("test-gateway", "t", "Użyj bramki testowej", func(s string) error {
		env = environmentPkg.Test
		return nil
	})

	RootCommand.AddCommand(authorization.AuthCommand)
	RootCommand.AddCommand(certificates.CertificatesCommand)
	RootCommand.AddCommand(syncInvoicesCommand)
	RootCommand.AddCommand(uploadCommand)
	RootCommand.CompletionOptions = cobra.CompletionOptions{DisableDefaultCmd: true}
}

func setContext(cmd *cobra.Command, _ []string) error {
	var err error

	if configFile != "" {
		if err = config.ReadConfig(configFile); err != nil {
			return err
		}
	}

	if err = logging.InitLogging(logOutput); err != nil {
		fmt.Printf("[ ERROR ] Błąd inicjalizacji logowania: %v", err)
		return err
	}

	logging.SeiLogger.Info("start programu")
	logging.SeiLogger.Info("wybrane środowisko", "env", env)

	var ctx = context.WithValue(
		cmd.Context(),
		appCtx.KeyEnvironment,
		env,
	)

	cmd.SetContext(ctx)

	return nil
}
