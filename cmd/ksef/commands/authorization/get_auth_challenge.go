package authorization

import (
	"fmt"

	"ksef/cmd/ksef/flags"
	v2 "ksef/internal/client/v2"
	"ksef/internal/client/v2/auth/token"
	"ksef/internal/resources"
	"ksef/internal/runtime"

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
	ctx := cmd.Context()
	vip := viper.GetViper()

	certsDB, err := resources.CertsDB(vip)
	if err != nil {
		return err
	}

	keyring, err := resources.Keyring(vip)
	if err != nil {
		return err
	}
	defer resources.Close(ctx, nil)

	output, err := cmd.Flags().GetString(flagOutput)
	if output == "" || err != nil {
		return fmt.Errorf("nie podano pliku wyjścia")
	}

	nip, err := runtime.GetNIP(vip)
	if err != nil {
		return err
	}

	authValidator := token.NewAuthHandler(
		token.WithNIP(nip),
		token.WithDumpChallenge(output),
		token.WithCertsDB(certsDB),
		token.WithKeyring(keyring),
	)

	rateLimitsState, err := resources.RateLimitsState(vip, nip)
	if err != nil {
		return err
	}

	cli, err := v2.NewClient(
		ctx,
		vip,
		v2.WithAuthValidator(authValidator),
		v2.WithKeyring(keyring),
		v2.WithRateLimitsStateStorage(rateLimitsState),
		v2.WithoutRestoringTokens(),
	)
	if err != nil {
		return err
	}

	return cli.AwaitTokenManagerLoop()
}
