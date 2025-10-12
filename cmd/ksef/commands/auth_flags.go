package commands

import (
	"errors"
	"ksef/cmd/ksef/flags"
	"ksef/internal/config"
	environmentPkg "ksef/internal/environment"
	kseftoken "ksef/internal/sei/api/client/v2/auth/ksef_token"
	"ksef/internal/sei/api/client/v2/auth/validator"
	"ksef/internal/sei/api/client/v2/auth/xades"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var (
	errUnableToSelectAuthBackend = errors.New("Nie udało się wybrać metody autoryzacji")
)

func authFlags(flagSet *pflag.FlagSet, cmd *cobra.Command) {
	flagSet.String(flags.FlagNameCertPath, "", "ścieżka do certyfikatu używanego do autoryzacji")
	flagSet.String(flags.FlagNameKSeFToken, "", "token KSeF lub nazwa zmiennej srodowiskowej która go zawiera")

	cmd.MarkFlagsOneRequired(flags.FlagNameCertPath, flags.FlagNameKSeFToken)
}

func authChallengeValidatorInstance(flagSet *pflag.FlagSet, nip string, env environmentPkg.Environment) (validator.AuthChallengeValidator, error) {
	var certPath string
	var err error

	if _, err = flagSet.GetString(flags.FlagNameKSeFToken); err != nil {
		return nil, err
	}
	if certPath, err = flagSet.GetString(flags.FlagNameCertPath); err != nil {
		return nil, err
	}

	apiConfig := config.GetConfig().APIConfig(env)

	if certPath != "" {
		// cert-based authentication
		return xades.NewAuthHandler(
			apiConfig,
			nip,
			certPath,
		), nil
	} else {
		// token-based authentication
		return kseftoken.NewKsefTokenHandler(
			apiConfig,
			nip,
		), nil
	}

	return nil, errUnableToSelectAuthBackend
}
