package xades

import (
	"errors"
	"ksef/internal/config"
	"ksef/internal/environment"
	v2 "ksef/internal/sei/api/client/v2"
	"ksef/internal/sei/api/client/v2/auth/validator"
	"ksef/internal/sei/api/client/v2/auth/xades"

	"github.com/spf13/cobra"
)

var debugCommand = &cobra.Command{
	Use:   "debug",
	Short: "inicjuje proces autoryzacji celem sprawdzenia poprawności działania certyfikatu",
	RunE:  authSessionDebug,
}

var (
	certFile      string
	challengeFile string
	signedFile    string
)

func init() {
	debugCommand.Flags().StringVarP(&certFile, "cert", "", "", "ścieżka do pliku certyfikatu służącego do podpisania wyzwania")
	debugCommand.Flags().StringVarP(&signedFile, "signed", "s", "", "ścieżka do *PODPISANEGO* pliku wyzwania")
	debugCommand.MarkFlagsOneRequired("cert", "signed")
	XadesCommand.AddCommand(debugCommand)
}

func authSessionDebug(cmd *cobra.Command, _ []string) error {
	var cfg = config.GetConfig()
	var env = environment.FromContext(cmd.Context())
	var authValidator validator.AuthChallengeValidator

	// there are couple of modes here:
	// 1. if the user passed path to a signed file - let's use it
	if signedFile != "" {
		authValidator = xades.NewSignedRequestHandler(
			cfg.APIConfig(env),
			signedFile,
		)
	} else if certFile != "" {
		// 2. if the user passed path to the certificate - we will do everything automatically
		nip, err := cmd.Flags().GetString("nip")
		if nip == "" || err != nil {
			return errors.New("brak numeru NIP")
		}
		authValidator = xades.NewAuthHandler(
			cfg.APIConfig(env),
			nip,
			certFile,
		)
	}

	cli, err := v2.NewClient(cmd.Context(), cfg, env, v2.WithAuthValidator(authValidator))
	if err != nil {
		return err
	}
	if err := cli.ObtainToken(); err != nil {
		return err
	}

	return cli.Logout()
}
