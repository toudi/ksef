package xades

import (
	"ksef/internal/config"
	"ksef/internal/environment"
	v2 "ksef/internal/sei/api/client/v2"
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
	debugCommand.Flags().StringVarP(&challengeFile, "challenge", "f", "", "ścieżka do pliku wyzwania (w przypadku podpisywania certyfikatem)")
	debugCommand.Flags().StringVarP(&signedFile, "signed", "s", "", "ścieżka do *PODPISANEGO* pliku wyzwania (w przypadku używania profilu zaufanego)")
	XadesCommand.AddCommand(debugCommand)
}

func authSessionDebug(cmd *cobra.Command, _ []string) error {
	var cfg = config.GetConfig()
	var env = environment.FromContext(cmd.Context())

	// there are couple of modes here:
	// 1. if the user passed path to a signed file - let's use it

	authValidator := xades.NewSignedRequestHandler(
		cfg.APIConfig(env),
		signedFile,
	)
	cli, err := v2.NewClient(cmd.Context(), cfg, env, v2.WithAuthValidator(authValidator))
	if err != nil {
		return err
	}
	if err := cli.ObtainToken(); err != nil {
		return err
	}

	return cli.Logout()

	return nil
}
