package authorization

import (
	"fmt"
	"ksef/cmd/ksef/flags"
	v2 "ksef/internal/client/v2"
	"ksef/internal/client/v2/auth/token"
	"ksef/internal/config"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	flagOutput = "output"
)

var initAuthCommand = &cobra.Command{
	Use:   "init",
	Short: "inicjalizuje sesję autoryzacyjną (pobiera wyzwanie i zapisuje do pliku)",
	RunE:  dumpAuthChallenge,
}

func init() {
	initAuthCommand.Flags().StringP(flagOutput, "o", "AuthTokenRequest.xml", "plik wyjściowy")
	flags.NIP(initAuthCommand.Flags())
	initAuthCommand.MarkFlagRequired(flags.FlagNameNIP)
	AuthCommand.AddCommand(initAuthCommand)
}

func dumpAuthChallenge(cmd *cobra.Command, _ []string) error {
	vip := viper.GetViper()
	nip, err := config.GetNIP(vip)
	if err != nil {
		return err
	}
	output, err := cmd.Flags().GetString(flagOutput)
	if output == "" || err != nil {
		return fmt.Errorf("nie podano pliku wyjścia")
	}

	authValidator := token.NewAuthHandler(
		config.GetGateway(vip),
		nip,
		token.WithDumpChallenge(output),
	)
	cli, err := v2.NewClient(cmd.Context(), config.GetGateway(vip), v2.WithAuthValidator(authValidator))
	if err != nil {
		return err
	}
	return cli.WaitForTokenManagerLoop()
}
