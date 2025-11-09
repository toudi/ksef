package authorization

import (
	"errors"
	"ksef/cmd/ksef/commands/authorization/challenge"
	"ksef/internal/certsdb"
	v2 "ksef/internal/client/v2"
	"ksef/internal/client/v2/auth/token"
	"ksef/internal/client/v2/auth/validator"
	"ksef/internal/config"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var debugCommand = &cobra.Command{
	Use:   "debug",
	Short: "inicjuje proces autoryzacji celem sprawdzenia poprawności działania certyfikatu",
	RunE:  authSessionDebug,
}

var (
	useCert       bool
	challengeFile string
	signedFile    string
)

func init() {
	debugCommand.Flags().BoolVar(&useCert, "cert", false, "spróbuj użyć certyfikatu")
	debugCommand.Flags().StringVarP(&signedFile, "signed", "s", "", "ścieżka do *PODPISANEGO* pliku wyzwania")
	debugCommand.MarkFlagsOneRequired("cert", "signed")
	AuthCommand.AddCommand(debugCommand)
}

func authSessionDebug(cmd *cobra.Command, _ []string) error {
	var authValidator validator.AuthChallengeValidator
	vip := viper.GetViper()
	var nip string
	var gateway = config.GetGateway(vip)
	var err error

	// there are couple of modes here:
	// 1. if the user passed path to a signed file - let's use it
	var initializerFuncs []func(handler *token.TokenHandler)

	if signedFile != "" {
		_, nip, err = challenge.GetNIPFromChallengeFile(signedFile)
		if err != nil {
			return err
		}
		initializerFuncs = append(initializerFuncs, token.WithSignedChallengeFile(signedFile))
	} else if useCert {
		// 2. if the user passed path to the certificate - we will do everything automatically
		nip, err = cmd.Flags().GetString("nip")
		if nip == "" || err != nil {
			return errors.New("brak numeru NIP")
		}
		// pick up cert from the database
		certsDB, err := certsdb.OpenOrCreate(gateway)
		if err != nil {
			return err
		}
		initializerFuncs = append(initializerFuncs, token.WithCertsDB(certsDB))
	}

	authValidator = token.NewAuthHandler(gateway, nip, initializerFuncs...)

	cli, err := v2.NewClient(cmd.Context(), gateway, v2.WithAuthValidator(authValidator))
	if err != nil {
		return err
	}
	if err := cli.ObtainToken(); err != nil {
		return err
	}

	return cli.Logout()
}
