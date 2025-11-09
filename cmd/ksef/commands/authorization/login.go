package authorization

import (
	"ksef/cmd/ksef/commands/authorization/challenge"
	v2 "ksef/internal/client/v2"
	"ksef/internal/client/v2/auth/token"
	"ksef/internal/config"

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

	_, nip, err := challenge.GetNIPFromChallengeFile(signedChallengeFile)
	if err != nil {
		return err
	}

	var gateway = config.GetGateway(viper.GetViper())

	var authValidator = token.NewAuthHandler(
		gateway,
		nip,
		token.WithSignedChallengeFile(signedChallengeFile),
	)

	cli, err := v2.NewClient(cmd.Context(), gateway, v2.WithAuthValidator(authValidator))
	if err != nil {
		return err
	}
	if err := cli.ObtainToken(); err != nil {
		return err
	}

	return nil
}
