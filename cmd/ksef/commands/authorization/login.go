package authorization

import (
	"ksef/cmd/ksef/commands/authorization/challenge"
	v2 "ksef/internal/client/v2"
	"ksef/internal/client/v2/auth/token"
	"ksef/internal/resources"
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
	ctx := cmd.Context()
	signedChallengeFile := args[0]
	vip := viper.GetViper()

	_, nip, err := challenge.GetNIPFromChallengeFile(signedChallengeFile)
	if err != nil {
		return err
	}
	runtime.SetNIP(vip, nip)

	keyring, err := resources.Keyring(vip)
	if err != nil {
		return err
	}
	defer resources.Close(ctx, nil)

	certsDB, err := resources.CertsDB(vip)
	if err != nil {
		return err
	}

	authValidator := token.NewAuthHandler(
		token.WithNIP(nip),
		token.WithSignedChallengeFile(signedChallengeFile),
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

	if err = cli.ObtainToken(); err != nil {
		return err
	}
	return cli.WaitForTokenManagerLoop()
}
