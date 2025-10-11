package authorization

import (
	"ksef/internal/config"
	"ksef/internal/environment"
	v2 "ksef/internal/sei/api/client/v2"
	"ksef/internal/sei/api/client/v2/auth/xades"

	"github.com/spf13/cobra"
)

var initAuthCommand = &cobra.Command{
	Use:   "init",
	Short: "inicjalizuje sesję autoryzacyjną (pobiera wyzwanie i zapisuje do pliku)",
	RunE:  dumpAuthChallenge,
}

var outputFile string

func init() {
	initAuthCommand.Flags().StringVarP(&outputFile, "output", "o", "AuthTokenRequest.xml", "plik wyjściowy")
	AuthCommand.AddCommand(initAuthCommand)
}

func dumpAuthChallenge(cmd *cobra.Command, _ []string) error {
	var cfg = config.GetConfig()
	var env = environment.FromContext(cmd.Context())

	authValidator := xades.NewChallengeDumperHandler(
		cfg.APIConfig(env),
		nip,
		outputFile,
	)
	cli, err := v2.NewClient(cmd.Context(), cfg, env, v2.WithAuthValidator(authValidator))
	if err != nil {
		return err
	}
	return cli.WaitForTokenManagerLoop()

}
