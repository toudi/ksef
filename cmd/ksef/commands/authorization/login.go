package authorization

import (
	"ksef/cmd/ksef/commands/authorization/challenge"
	"ksef/internal/certsdb"
	v2 "ksef/internal/client/v2"
	"ksef/internal/client/v2/auth"
	"ksef/internal/client/v2/auth/token"
	"ksef/internal/runtime"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var loginCommand = &cobra.Command{
	Use:   "login AuthTokenRequest.signed.xml",
	Short: "używa podpisanego pliku wyzwania aby pobrać token sesyjny",
	Args:  cobra.ExactArgs(1),
	RunE:  login,
}

func init() {
	loginCommand.Flags().SortFlags = false

	AuthCommand.AddCommand(loginCommand)
}

func login(cmd *cobra.Command, args []string) error {
	var signedChallengeFile = args[0]
	vip := viper.GetViper()
	vip.Set(auth.FlagDoNotRestoreTokens, "true")
	vip.Set(auth.FlagExitAfterPersistingToken, "true")

	_, nip, err := challenge.GetNIPFromChallengeFile(signedChallengeFile)
	if err != nil {
		return err
	}
	runtime.SetNIP(vip, nip)

	env := runtime.GetGateway(viper.GetViper())
	certsDB, err := certsdb.OpenOrCreate(env)
	if err != nil {
		return err
	}

	var authValidator = token.NewAuthHandler(
		vip,
		token.WithSignedChallengeFile(signedChallengeFile),
		token.WithCertsDB(certsDB),
	)

	cli, err := v2.NewClient(cmd.Context(), vip, v2.WithAuthValidator(authValidator))
	if err != nil {
		return err
	}

	if err = cli.ObtainToken(); err != nil {
		return err
	}
	return cli.WaitForTokenManagerLoop()

}
