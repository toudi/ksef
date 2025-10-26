package commands

import (
	"context"
	"fmt"
	"ksef/cmd/ksef/commands/authorization"
	"ksef/cmd/ksef/commands/certificates"
	appCtx "ksef/cmd/ksef/context"
	environmentPkg "ksef/internal/environment"
	"ksef/internal/logging"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	RootCommand.PersistentFlags().StringVarP(&configFile, "config", "c", "config.yaml", "lokalizacja pliku konfiguracyjnego")
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
		if err = viper.ReadInConfig(); err != nil {
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
