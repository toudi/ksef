package commands

import (
	"context"
	"ksef/cmd/ksef/commands/authorization"
	appCtx "ksef/cmd/ksef/context"
	"ksef/internal/config"
	environmentPkg "ksef/internal/environment"

	"github.com/spf13/cobra"
)

var RootCommand = &cobra.Command{
	Short:             "ksef",
	Long:              "KSeF - aplikacja kliencka",
	PersistentPreRunE: setContext,
}

var (
	env        environmentPkg.Environment = environmentPkg.Production
	configFile string
)

func init() {
	RootCommand.PersistentFlags().StringVarP(&configFile, "config", "c", "", "lokalizacja pliku konfiguracyjnego")
	RootCommand.PersistentFlags().BoolFuncP("test-gateway", "t", "Użyj bramki testowej", func(s string) error {
		env = environmentPkg.Test
		return nil
	})

	RootCommand.AddCommand(authorization.AuthCommand)
	RootCommand.CompletionOptions = cobra.CompletionOptions{DisableDefaultCmd: true}
}

func setContext(cmd *cobra.Command, _ []string) error {
	if configFile != "" {
		if err := config.ReadConfig(configFile); err != nil {
			return err
		}
	}

	var ctx = context.WithValue(
		cmd.Context(),
		appCtx.KeyEnvironment,
		env,
	)

	cmd.SetContext(ctx)

	return nil
}
