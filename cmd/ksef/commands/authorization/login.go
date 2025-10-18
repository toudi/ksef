package authorization

import (
	"ksef/cmd/ksef/flags"
	"ksef/internal/config"
	"ksef/internal/environment"
	v2 "ksef/internal/sei/api/client/v2"
	"ksef/internal/sei/api/client/v2/auth/validator"
	"ksef/internal/sei/api/client/v2/auth/xades"

	"github.com/spf13/cobra"
)

var loginCommand = &cobra.Command{
	Use:   "login AuthTokenRequest.signed.xml",
	Short: "używa podpisanego pliku wyzwania aby pobrać token sesyjny",
	Args:  cobra.ExactArgs(1),
	RunE:  login,
}

func init() {
	loginCommand.Flags().SortFlags = false
	loginCommand.MarkFlagRequired(flags.FlagNameNIP)

	AuthCommand.AddCommand(loginCommand)
}

func login(cmd *cobra.Command, args []string) error {
	var signedChallengeFile = args[0]

	var cfg = config.GetConfig()
	var env = environment.FromContext(cmd.Context())
	var authValidator validator.AuthChallengeValidator = xades.NewSignedRequestHandler(
		cfg.APIConfig(env),
		signedChallengeFile,
	)
	nip, err := cmd.Flags().GetString(flags.FlagNameNIP)
	if err != nil {
		return err
	}

	cli, err := v2.NewClient(cmd.Context(), cfg, env, v2.WithAuthValidator(authValidator))
	if err != nil {
		return err
	}
	if err := cli.ObtainToken(); err != nil {
		return err
	}

	if err = cli.PersistTokens(env, nip); err != nil {
		return err
	}

	return nil
}
